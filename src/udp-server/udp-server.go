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

	"github.com/segmentio/kafka-go"
)

func handleClientMessage(clientAddress *net.UDPAddr, message []byte, kafkaProducer *kafka.Writer,
	kafkaTimeout time.Duration) {
	// Parse packet header
	header := PacketHeader{}
	err := binary.Read(bytes.NewReader(message), binary.LittleEndian, &header)
	if err != nil {
		log.Println(clientAddress, "- Error parsing header:", err)
		return
	}

	headerSize := binary.Size(header)
	messageNoHeader := message[headerSize:]

	_, prst := PACKET_MAP[header.M_packetId]
	if !prst {
		log.Println(clientAddress, "- Unknown packet ID:", header.M_packetId)
		return
	}

	// Generic packet variable and topic name
	var packet any = PACKET_MAP[header.M_packetId]
	var topicName string = PACKET_TOPIC_MAP[header.M_packetId]

	var packetData any

	switch packet.(type) {
	case PacketMotionData:
		packetData, err = parsePacketData(messageNoHeader, &PacketMotionData{})
	case PacketSessionData:
		packetData, err = parsePacketData(messageNoHeader, &PacketSessionData{})
	case PacketLapData:
		packetData, err = parsePacketData(messageNoHeader, &PacketLapData{})
	case GenericEvent:
		eventCode := string(messageNoHeader[:4])
		_, prst = EVENT_MAP[eventCode]
		if !prst {
			log.Println(clientAddress, "- Unknown event code:", eventCode)
			return
		}
		packet = EVENT_MAP[eventCode]

		switch packet.(type) {
		case PacketEventSSTA:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventSSTA{})
		case PacketEventSEND:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventSEND{})
		case PacketEventFTLP:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventFTLP{})
		case PacketEventRTMT:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventRTMT{})
		case PacketEventDRSE:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventDRSE{})
		case PacketEventDRSD:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventDRSD{})
		case PacketEventTMPT:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventTMPT{})
		case PacketEventCHQF:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventCHQF{})
		case PacketEventRCWN:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventRCWN{})
		case PacketEventPENA:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventPENA{})
		case PacketEventSPTP:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventSPTP{})
		case PacketEventSTLG:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventSTLG{})
		case PacketEventDTSV:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventDTSV{})
		case PacketEventSGSV:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventSGSV{})
		case PacketEventFLBK:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventFLBK{})
		case PacketEventBUTN:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventBUTN{})
		case PacketEventRDFL:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventRDFL{})
		case PacketEventOVTK:
			packetData, err = parsePacketData(messageNoHeader, &PacketEventOVTK{})
		default:
			log.Println(clientAddress, "- Unknown event code:", eventCode)
			return
		}
		topicName = TOPIC_EVENT_DATA
	case PacketParticipantsData:
		packetData, err = parsePacketData(messageNoHeader, &PacketParticipantsData{})
	case PacketCarSetupData:
		packetData, err = parsePacketData(messageNoHeader, &PacketCarSetupData{})
	case PacketCarTelemetryData:
		packetData, err = parsePacketData(messageNoHeader, &PacketCarTelemetryData{})
	case PacketCarStatusData:
		packetData, err = parsePacketData(messageNoHeader, &PacketCarStatusData{})
	case PacketFinalClassificationData:
		packetData, err = parsePacketData(messageNoHeader, &PacketFinalClassificationData{})
	case PacketLobbyInfoData:
		packetData, err = parsePacketData(messageNoHeader, &PacketLobbyInfoData{})
	case PacketCarDamageData:
		packetData, err = parsePacketData(messageNoHeader, &PacketCarDamageData{})
	case PacketSessionHistoryData:
		packetData, err = parsePacketData(messageNoHeader, &PacketSessionHistoryData{})
	case PacketTyreSetsData:
		packetData, err = parsePacketData(messageNoHeader, &PacketTyreSetsData{})
	case PacketMotionExData:
		packetData, err = parsePacketData(messageNoHeader, &PacketMotionExData{})
	default:
		log.Println(clientAddress, "- Unknown packet type:", packet)
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

	sendMessageToKafka(kafkaProducer, topicName, string(jsonData), kafkaTimeout)
}

func parsePacketData(messageNoHeader []byte, packet any) (any, error) {
	err := binary.Read(bytes.NewReader(messageNoHeader), binary.LittleEndian, packet)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func createKafkaTopics(address string, port string, topics []string) error {
	conn, err := kafka.Dial("tcp", address+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	return controllerConn.CreateTopics(topicConfigs...)
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
		log.Println("No PORT environment variable set, using default port", PORT)
		srv_port = PORT
	}
	kafka_address := os.Getenv("KAFKA_ADDRESS")
	if kafka_address == "" {
		log.Println("No KAFKA_ADDRESS environment variable set, using default address", KAFKA_ADDRESS)
		kafka_address = KAFKA_ADDRESS
	}
	kafka_port := os.Getenv("KAFKA_PORT")
	if kafka_port == "" {
		log.Println("No KAFKA_PORT environment variable set, using default port", KAFKA_PORT)
		kafka_port = KAFKA_PORT
	}

	// Setup Apache Kafka topics
	err := createKafkaTopics(kafka_address, kafka_port, MESSAGE_TOPICS[:])
	if err != nil {
		log.Fatalln("Error creating Kafka topics:", err)
	}
	log.Println("Kafka topics created successfully:", MESSAGE_TOPICS)

	// Setup Kafka producer
	producer := &kafka.Writer{
		Addr:     kafka.TCP(kafka_address + ":" + kafka_port),
		Balancer: &kafka.LeastBytes{},
	}
	defer producer.Close()

	// Setup UDP server
	addr, err := net.ResolveUDPAddr("udp", ADDR+":"+srv_port)
	if err != nil {
		log.Fatalln(err)
	}

	con, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	log.Println("UDP server listening on", ADDR+":"+srv_port)

	// Endless receive loop
	for {
		buffer := make([]byte, MAX_BUFFER_SIZE)
		n, clientAddr, err := con.ReadFromUDP(buffer)
		if err != nil {
			log.Println(clientAddr, "- Error reading:", err)
		}

		go handleClientMessage(clientAddr, buffer[:n], producer, MESSAGE_TIMEOUT)
	}
}
