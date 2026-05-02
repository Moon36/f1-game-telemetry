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
  - A map of string keys and string values containing the contents of the CSV file, or an error if the file could not be
    read.
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

/*
Updates the Redis database with participant data from the given PacketParticipantsData struct. Each participant's data
is stored in Redis under a key formatted as "participant:{index}", where {index} is the participant's index in the
M_participants slice of the PacketParticipantsData struct.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - participantsPacket: The PacketParticipantsData struct containing the participant data to be stored in Redis.
*/
func updateRedisWithParticipantData(ctx context.Context,
	redisClient *redis.Client,
	participantsPacket packets.PacketParticipantsData) {
	wg := sync.WaitGroup{}
	errorChan := make(chan error, len(participantsPacket.M_participants))

	for i, participant := range participantsPacket.M_participants {
		wg.Add(1)
		go translateAndStoreParticipantData(ctx, redisClient, i, participant, errorChan, &wg)
	}

	wg.Wait()
	close(errorChan)

	if len(errorChan) != 0 {
		log.Println("some participants failed to be updated in Redis. Total errors:",
			fmt.Sprintf("%d/%d:", len(errorChan), len(participantsPacket.M_participants)))

		for err := range errorChan {
			log.Printf("an error occurred while updating Redis with participant data: %v", err)
		}
		return
	}

	log.Println("successfully updated Redis with new participant data.")
}

func translateAndStoreParticipantData(ctx context.Context,
	redisClient *redis.Client,
	idx int,
	participant packets.ParticipantData,
	errorChan chan error,
	wg *sync.WaitGroup) {
	defer wg.Done()

	resolvedParticipant, err := resolveParticipantIDs(ctx, redisClient, participant)
	if err != nil {
		errorChan <- err
		return
	}

	redisKey := fmt.Sprintf("participant:%d", idx)
	err = redisClient.JSONSet(ctx, redisKey, "$", resolvedParticipant).Err()
	if err != nil {
		errorChan <- err
	}
}

/*
Resolves the IDs in the given ParticipantData struct to their corresponding names using the data stored in Redis.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - participant: The ParticipantData struct containing the participant data with IDs to be resolved.
*/
func resolveParticipantIDs(ctx context.Context,
	redisClient *redis.Client,
	participant packets.ParticipantData) (packets.StoreParticipant, error) {
	storeParticipant := packets.StoreParticipant{
		M_aiControlled:    participant.M_aiControlled,
		M_networkId:       participant.M_networkId,
		M_myTeam:          participant.M_myTeam,
		M_raceNumber:      participant.M_raceNumber,
		M_name:            participant.M_name,
		M_yourTelemetry:   participant.M_yourTelemetry,
		M_showOnlineNames: participant.M_showOnlineNames,
	}
	err := error(nil)
	// Driver ID to name
	driverName := "Player"
	if participant.M_driverId != 255 {
		driverName, err = getNameById(ctx, redisClient, "csv:drivers", participant.M_driverId)
		if err != nil {
			return packets.StoreParticipant{}, err
		}
	}
	copy(storeParticipant.M_driverName[:], []byte(driverName))

	// Team ID to name
	teamName, err := getNameById(ctx, redisClient, "csv:teams", participant.M_teamId)
	if err != nil {
		return packets.StoreParticipant{}, err
	}
	copy(storeParticipant.M_teamName[:], []byte(teamName))

	// Nationality ID to name
	nationality, err := getNameById(ctx, redisClient, "csv:nationalities", participant.M_nationality)
	if err != nil {
		return packets.StoreParticipant{}, err
	}
	copy(storeParticipant.M_teamName[:], []byte(nationality))

	// Platform ID to name
	platformName, err := getNameById(ctx, redisClient, "csv:platforms", participant.M_platform)
	if err != nil {
		return packets.StoreParticipant{}, err
	}
	copy(storeParticipant.M_teamName[:], []byte(platformName))

	return storeParticipant, nil
}

/*
Looks for the given ID in the given hash-set in Redis and returns the associated value.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - hashSet: The hash-set/key to use for the lookup.
  - id: The field to look up.
*/
func getNameById(ctx context.Context, redisClient *redis.Client, hashSet string, id uint8) (string, error) {
	value, err := redisClient.HGet(ctx, hashSet, fmt.Sprintf("%d", id)).Result()
	if err != nil {
		return "", err
	}
	return value, nil
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

		log.Println("Received participants data.")
		participantsPacket := packets.PacketParticipantsData{}
		err = json.Unmarshal(msg.Value, &participantsPacket)
		if err != nil {
			log.Printf("failed to unmarshal participant data: %v", err)
			continue
		}

		log.Println("Refreshing Redis with new data...")
		updateRedisWithParticipantData(ctx, rdb, participantsPacket)
	}

	if err := consumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
