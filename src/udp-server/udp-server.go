package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	common "github.com/moon36/f1-game-telemetry/src/internal"
	"github.com/moon36/f1-game-telemetry/src/internal/mappers"
	"github.com/moon36/f1-game-telemetry/src/internal/packets"

	"github.com/segmentio/kafka-go"
)

/*
Handles a client message. It parses the header, maps it to a version specific struct representation and submits the
parsed and marshalled message with the corresponding topic to the message queue.

Parameters:
  - clientAddress: The address of the client that sent the message.
  - message: The raw bytes of the message received from the client.
  - kafkaProducer: A Kafka producer instance used to send messages to a Kafka topic.
  - kafkaTimeout: The timeout duration for sending messages through the Kafka producer.
*/
func handleClientMessage(clientAddress *net.UDPAddr, message []byte, kafkaProducer *kafka.Writer,
	kafkaTimeout time.Duration) {
	// Parse packet header
	header := packets.PacketHeader{}
	err := binary.Read(bytes.NewReader(message), binary.LittleEndian, &header)
	if err != nil {
		log.Println(clientAddress, "- Error parsing header:", err)
		return
	}

	mapper, prst := mappers.PACKET_MAPPER_MAP[header.M_packetFormat]
	if !prst {
		log.Println(clientAddress, "- No packet mapper registered (yet) for packet format:", header.M_packetFormat)
		return
	}

	packetData, topicName, err := mapper.MapPacket(header, message)

	if err != nil {
		log.Println(clientAddress, "- Could not map packet of packet format:", header.M_packetFormat, "\n",
			"Failed with error", err)
		return
	}

	err = parsePacketData(message, packetData)
	if err != nil {
		log.Println(clientAddress, "- Error parsing packet data:", err)
		return
	}

	// Marshal data
	jsonData, err := json.Marshal(packetData)
	if err != nil {
		log.Println(clientAddress, "- Error marshaling to JSON:", err)
		return
	}

	err = sendMessageToKafka(kafkaProducer, topicName, string(jsonData), kafkaTimeout)
	if err != nil {
		log.Println(clientAddress, "- Error sending message to Kafka:", err)
		return
	}
}

func parsePacketData(message []byte, packet any) error {
	err := binary.Read(bytes.NewReader(message), binary.LittleEndian, packet)
	if err != nil {
		return err
	}
	return nil
}

/*
Creates Kafka topics using the provided address, port, and topic list.

Parameters:
  - address: The address of the Kafka broker.
  - port: The port number of the Kafka broker.
  - topics: A list of topic names to create.

Returns:
  - error: An error if any step fails, otherwise nil.
*/
func createKafkaTopics(address string, port string, topics []string) error {
	conn, err := kafka.Dial("tcp", address+":"+port)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
	}()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer func() {
		err = controllerConn.Close()
	}()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	err = controllerConn.CreateTopics(topicConfigs...)

	return err
}

/*
Submits a JSON message to a Kafka topic with the given timeout.
If the message cannot be sent within the timeout, it will return an error.

Parameters:
  - kafkaProducer: A Kafka producer instance.
  - topic: The topic to send the message to.
  - jsonMessage: The JSON message to be sent.
  - timeout: The maximum time to wait for the message to be sent.

Returns:
  - error: If an error occurs during the message submission.
*/
func sendMessageToKafka(kafkaProducer *kafka.Writer, topic string, jsonMessage string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := kafkaProducer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Value: []byte(jsonMessage),
		},
	)
	select {
	case <-ctx.Done():
		log.Println("Kafka message timed out:", ctx.Err())
		return ctx.Err()
	default:
		if err != nil {
			log.Println("Error sending message to Kafka:", err)
			return err
		}
	}

	return nil
}

func main() {
	// Get environment variables
	srv_port := os.Getenv("PORT")
	if srv_port == "" {
		log.Println("No PORT environment variable set, using default port", common.PORT)
		srv_port = common.PORT
	}
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

	// Setup Apache Kafka topics
	err := createKafkaTopics(kafka_address, kafka_port, packets.MESSAGE_TOPICS[:])
	if err != nil {
		log.Fatalln("Error creating Kafka topics:", err)
		return
	}
	log.Println("Kafka topics created successfully:", packets.MESSAGE_TOPICS)

	// Setup Kafka producer
	producer := &kafka.Writer{
		Addr:     kafka.TCP(kafka_address + ":" + kafka_port),
		Balancer: &kafka.LeastBytes{},
	}
	defer func() {
		err = producer.Close()
		if err != nil {
			log.Println("Error closing Kafka producer:", err)
		}
	}()

	// Setup UDP server
	addr, err := net.ResolveUDPAddr("udp", common.ADDR+":"+srv_port)
	if err != nil {
		log.Fatalln(err)
		return
	}

	con, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer func() {
		err = con.Close()
		if err != nil {
			log.Println("Error closing UDP connection:", err)
		}
	}()

	log.Println("UDP server listening on", common.ADDR+":"+srv_port)

	// Endless receive loop
	for {
		buffer := make([]byte, packets.MAX_BUFFER_SIZE)
		n, clientAddr, err := con.ReadFromUDP(buffer)
		if err != nil {
			log.Println(clientAddr, "- Error reading:", err)
		}

		go handleClientMessage(clientAddr, buffer[:n], producer, common.MESSAGE_TIMEOUT)
	}
}
