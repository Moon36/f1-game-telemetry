package main

import "time"

// Server constants
const ADDR = "0.0.0.0"
const PORT = "8888"

// Kafka constants
const KAFKA_ADDRESS = "localhost" // Use Docker service name
const KAFKA_PORT = "9092"
const MESSAGE_TIMEOUT = 5 * time.Second
