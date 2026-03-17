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

	common "github.com/moon36/f1-game-telemetry/src/internal"
	"github.com/moon36/f1-game-telemetry/src/internal/packets"

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
  - filePath: The path to the CSV file to be read.

Returns:
  - A map of string keys and string values containing the contents of the CSV file, or an error if the file could not be read.
*/
func readCsvFile(fs embed.FS, filePath string) (map[string]string, error) {
	f, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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
  - csvFilePaths: A slice of file paths to the CSV files that contain the static data to be loaded into Redis.

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

func updateRedisWithParticipantData(ctx context.Context, redisClient *redis.Client, participantsPacket packets.PacketParticipantsData) {
	for i, participant := range participantsPacket.M_participants {
		redisKey := fmt.Sprintf("participant:%d", i)
		err := redisClient.JSONSet(ctx, redisKey, "$", participant).Err()
		if err != nil {
			log.Printf("failed to store participant data in Redis: %v\n", err)
		}
	}
	log.Println("Successfully updated Redis with new participant data.")
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
	fmt.Print("Setting up Redis client at address: ", redis_address, ":", redis_port, "\n")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_address + ":" + redis_port,
		Password: "",
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
		GroupID: "db-manager-consumer",
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

		if msg.Topic != packets.TOPIC_PARTICIPANT_DATA {
			continue
		}

		fmt.Println("Received participants data.")
		participantsPacket := packets.PacketParticipantsData{}
		err = json.Unmarshal(msg.Value, &participantsPacket)
		if err != nil {
			log.Printf("failed to unmarshal participant data: %v", err)
			continue
		}

		fmt.Println("Refreshing Redis with new data...")
		go updateRedisWithParticipantData(ctx, rdb, participantsPacket)
	}

	if err := consumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
