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

	"github.com/moon36/f1-game-telemetry/src/udp-server/packets"

	"github.com/segmentio/kafka-go"
)

func handleClientMessage(clientAddress *net.UDPAddr, message []byte, kafkaProducer *kafka.Writer,
	kafkaTimeout time.Duration) {
	// Parse packet header
	header := packets.PacketHeader{}
	err := binary.Read(bytes.NewReader(message), binary.LittleEndian, &header)
	if err != nil {
		log.Println(clientAddress, "- Error parsing header:", err)
		return
	}

	headerSize := binary.Size(header)
	messageNoHeader := message[headerSize:]

	_, prst := packets.PACKET_MAP[header.M_packetId]
	if !prst {
		log.Println(clientAddress, "- Unknown packet ID:", header.M_packetId)
		return
	}

	// Generic packet variable and topic name
	var packet any = packets.PACKET_MAP[header.M_packetId]
	var topicName string = packets.PACKET_TOPIC_MAP[header.M_packetId]

	var packetData any

	switch packet.(type) {
	case packets.PacketMotionData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketMotionData{})
	case packets.PacketSessionData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketSessionData{})
	case packets.PacketLapData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketLapData{})
	case packets.GenericEvent:
		eventCode := string(messageNoHeader[:4])
		_, prst = packets.EVENT_MAP[eventCode]
		if !prst {
			log.Println(clientAddress, "- Unknown event code:", eventCode)
			return
		}
		packet = packets.EVENT_MAP[eventCode]

		switch packet.(type) {
		case packets.PacketEventSSTA:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventSSTA{})
		case packets.PacketEventSEND:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventSEND{})
		case packets.PacketEventFTLP:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventFTLP{})
		case packets.PacketEventRTMT:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventRTMT{})
		case packets.PacketEventDRSE:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventDRSE{})
		case packets.PacketEventDRSD:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventDRSD{})
		case packets.PacketEventTMPT:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventTMPT{})
		case packets.PacketEventCHQF:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventCHQF{})
		case packets.PacketEventRCWN:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventRCWN{})
		case packets.PacketEventPENA:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventPENA{})
		case packets.PacketEventSPTP:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventSPTP{})
		case packets.PacketEventSTLG:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventSTLG{})
		case packets.PacketEventDTSV:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventDTSV{})
		case packets.PacketEventSGSV:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventSGSV{})
		case packets.PacketEventFLBK:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventFLBK{})
		case packets.PacketEventBUTN:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventBUTN{})
		case packets.PacketEventRDFL:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventRDFL{})
		case packets.PacketEventOVTK:
			packetData, err = parsePacketData(messageNoHeader, &packets.PacketEventOVTK{})
		default:
			log.Println(clientAddress, "- Unknown event code:", eventCode)
			return
		}
		topicName = packets.TOPIC_EVENT_DATA
	case packets.PacketParticipantsData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketParticipantsData{})
	case packets.PacketCarSetupData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketCarSetupData{})
	case packets.PacketCarTelemetryData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketCarTelemetryData{})
	case packets.PacketCarStatusData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketCarStatusData{})
	case packets.PacketFinalClassificationData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketFinalClassificationData{})
	case packets.PacketLobbyInfoData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketLobbyInfoData{})
	case packets.PacketCarDamageData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketCarDamageData{})
	case packets.PacketSessionHistoryData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketSessionHistoryData{})
	case packets.PacketTyreSetsData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketTyreSetsData{})
	case packets.PacketMotionExData:
		packetData, err = parsePacketData(messageNoHeader, &packets.PacketMotionExData{})
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
	err := createKafkaTopics(kafka_address, kafka_port, packets.MESSAGE_TOPICS[:])
	if err != nil {
		log.Fatalln("Error creating Kafka topics:", err)
	}
	log.Println("Kafka topics created successfully:", packets.MESSAGE_TOPICS)

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
		buffer := make([]byte, packets.MAX_BUFFER_SIZE)
		n, clientAddr, err := con.ReadFromUDP(buffer)
		if err != nil {
			log.Println(clientAddr, "- Error reading:", err)
		}

		go handleClientMessage(clientAddr, buffer[:n], producer, MESSAGE_TIMEOUT)
	}
}
