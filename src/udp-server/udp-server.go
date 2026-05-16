package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	common "github.com/moon36/f1-game-telemetry/src/internal"
	"github.com/moon36/f1-game-telemetry/src/internal/packets"

	"github.com/segmentio/kafka-go"
)

func processPacket23(packetID uint8, message []byte) (any, string, error) {
	createPacket, prst := packets.PACKET_MAP_23[packetID]
	if !prst {
		return nil, "", fmt.Errorf("unknown packet ID: %d", packetID)
	}

	var packetData any
	var err error
	var topicName string

	if packetID == packets.EVENT_DATA_ID {
		if len(message) < 4 {
			return nil, "", fmt.Errorf("supposed event data packet ('23) does not contain event code")
		}
		eventCode := string(message[:4])
		createPacket, prst = packets.EVENT_MAP_23[eventCode]
		if !prst {
			return nil, "", fmt.Errorf("unknown event code: %s", eventCode)
		}
	}
	packetData = createPacket()
	topicName = packets.PACKET_TOPIC_MAP_23[packetID]

	err = parsePacketData(message, packetData)
	return packetData, topicName, err
}

func handleClientMessage(clientAddress *net.UDPAddr, message []byte, kafkaProducer *kafka.Writer,
	kafkaTimeout time.Duration) {
	// Parse packet header
	header := packets.PacketHeader{}
	err := binary.Read(bytes.NewReader(message), binary.LittleEndian, &header)
	if err != nil {
		log.Println(clientAddress, "- Error parsing header:", err)
		return
	}

	// Generic packet variable and topic name
	var packetData any
	var topicName string

	switch header.M_packetFormat {
	case packets.PACKET_FORMAT_ID_23:
		packetData, topicName, err = processPacket23(header.M_packetId, message)
	default:
		log.Println(clientAddress, "- Unknown packet format ID:", header.M_packetFormat,
			"(This format might not be supported yet)")
		return
	}

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

	// TODO: Remove me
	log.Println("Message sent to Kafka topic:", topic)
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
