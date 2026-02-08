package common

import "time"

// Server constants
const ADDR = "0.0.0.0"
const PORT = "8888"

// Kafka constants
const KAFKA_ADDRESS = "localhost" // Use Docker service name
const KAFKA_PORT = "9092"
const MESSAGE_TIMEOUT = 5 * time.Second

// Redis constants
const REDIS_ADDRESS = "localhost" // Use Docker service name
const REDIS_PORT = "6379"
const REDIS_DYNAMIC_DATA_TTL = 10
