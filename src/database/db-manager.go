package main

import (
	"context"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	common "github.com/moon36/f1-game-telemetry/src/internal"
	"github.com/moon36/f1-game-telemetry/src/internal/packets"
	v23 "github.com/moon36/f1-game-telemetry/src/internal/packets/v23"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

//go:embed storage/*.csv
var csvFiles embed.FS

/*
Returns the file name without its extension from the given file name string.
From: https://gist.github.com/ivanzoid/129460aa08aff72862a534ebe0a9ae30?permalink_comment_id=3733302#gistcomment-3733302

Parameters:
  - fileName: The name of the file from which to extract the base name.
*/
func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

/*
Reads a CSV file from the specified file path and returns its contents as a map of string keys and string values.
Each row in the CSV file is expected to have two columns, where the first column is used as the key and the second
column as the value in the resulting map.

Parameters:
  - fs: The embedded file system containing the CSV file to be read.
  - filePath: The path to the CSV file to be read.

Returns:
  - A map of string keys and string values containing the contents of the CSV file, or an error if the file could not be
    read.
*/
func readCsvFile(fs embed.FS, filePath string) (map[string]string, error) {
	f, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("failed to close file %s: %v", filePath, err)
		}
	}()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)

	for _, record := range records {
		if len(record) != 2 {
			return nil, fmt.Errorf("invalid record in CSV file %s: %v", filePath, record)
		}
		if record[0] == "id" || record[0] == "bit_flag" {
			continue
		}
		entries[record[0]] = record[1]
	}
	return entries, nil
}

/*
Sets up the static data in Redis by reading the specified CSV files and storing their contents as hash maps.
Each CSV file is expected to have two columns, where the first column is used as the key and the second column as the
value in the Redis hash.
The Redis keys are prefixed with "csv:" followed by the base name of the CSV file.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - fs: The embedded file system containing the CSV files to be read.
  - folderPath: The path to the folder containing the CSV files.

Returns:
  - An error if any issues occur while reading the CSV files or storing the data in Redis, otherwise nil.
*/
func setupStaticData(ctx context.Context, redisClient *redis.Client, fs embed.FS, folderPath string) error {
	files, err := fs.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, filePath := range files {
		fpath := filepath.Join(folderPath, filePath.Name())

		fileData, err := readCsvFile(fs, fpath)
		if err != nil {
			return fmt.Errorf("failed to read CSV file %s: %v", fpath, err)
		}

		redisKey := fmt.Sprintf("csv:%s", fileNameWithoutExtension(filepath.Base(filePath.Name())))

		err = redisClient.HSet(ctx, redisKey, fileData).Err()
		if err != nil {
			return fmt.Errorf("failed to store data from file %s in Redis: %v", filePath.Name(), err)
		}
	}
	return nil
}

func main() {
	// Get environment variables
	kafka_address := os.Getenv("KAFKA_ADDRESS")
	if kafka_address == "" {
		log.Println("No KAFKA_ADDRESS environment variable set, using default address", common.KAFKA_ADDRESS)
		kafka_address = common.KAFKA_ADDRESS
	}
	kafka_port := os.Getenv("KAFKA_PORT")
	if kafka_port == "" {
		log.Println("No KAFKA_PORT environment variable set, using default port", common.KAFKA_PORT)
		kafka_port = common.KAFKA_PORT
	}
	redis_address := os.Getenv("REDIS_ADDRESS")
	if redis_address == "" {
		log.Println("No REDIS_ADDRESS environment variable set, using default address", common.REDIS_ADDRESS)
		redis_address = common.REDIS_ADDRESS
	}
	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		log.Println("No REDIS_PORT environment variable set, using default port", common.REDIS_PORT)
		redis_port = common.REDIS_PORT
	}

	// Setup Redis client
	var ctx = context.Background()
	log.Println("Setting up Redis client at address:", redis_address, " and port:", redis_port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_address + ":" + redis_port,
		Username: USERNAME,
		Password: PASSWORD,
		DB:       0,
	})
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Fatal("failed to close Redis client:", err)
		}
	}()

	// Load static data into Redis
	err := setupStaticData(ctx, rdb, csvFiles, "storage")
	if err != nil {
		log.Fatal("failed to setup static data:", err)
	}

	// Setup Kafka consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafka_address + ":" + kafka_port},
		Topic:   packets.TOPIC_PARTICIPANT_DATA,
		GroupID: KAFKA_CONSUMER_GROUP_ID,
	})
	defer func() {
		err := consumer.Close()
		if err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	log.Println("Database manager started, consuming from Kafka topic:", packets.TOPIC_PARTICIPANT_DATA)
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("an error occurred while trying to read message from Kafka:", err)
			break
		}

		log.Println("Received participants data.")

		var basePacket BasePacket
		err = json.Unmarshal(msg.Value, &basePacket)
		if err != nil {
			log.Printf("failed to unmarshal base packet: %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		var ch chan ChannelData

		switch basePacket.M_header.M_packetFormat {
		case v23.PACKET_FORMAT_ID:
			log.Println("Processing packet with format ID 2023")
			participantsPacket := v23.PacketParticipantsData{}
			err = json.Unmarshal(msg.Value, &participantsPacket)

			if err != nil {
				log.Printf("failed to unmarshal participant data: %v", err)
				continue
			}

			ch = make(chan ChannelData, len(participantsPacket.M_participants))

			for idx, participant := range participantsPacket.M_participants {
				wg.Add(1)
				go enrichParticipant23(ctx, rdb, uint8(idx), participant, ch)
			}
		default:
			log.Println("Received packet with unknown format ID:", basePacket.M_header.M_packetFormat)
			continue
		}

		go func() {
			// Wait for all goroutines to finish and then close the channel to signal that no more data will be sent
			wg.Wait()
			close(ch)
		}()

		// Consume enriched participant data from channel until closed and store in Redis
		for chData := range ch {
			err := chData.Error
			if err != nil {
				log.Printf("error enriching participant data for participant with ID %d: %v",
					chData.Idx, err)
				continue
			}
			redisKey := fmt.Sprintf("participant:%d", chData.Idx)
			err = rdb.JSONSet(ctx, redisKey, "$", chData.ParticipantData).Err()
			if err != nil {
				log.Printf("failed to store participant data in Redis for participant with ID %d: %v",
					chData.Idx, err)
			}
		}

		log.Println("Refreshing Redis with new data...")
	}

	if err := consumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
