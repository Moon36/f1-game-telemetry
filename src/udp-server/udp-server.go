package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"os"
)

func handleClientMessage(clientAddress *net.UDPAddr, message []byte) {
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

	// TODO: Send data to message broker
	if topicName == TOPIC_CAR_MOTION_DATA {
		log.Println(clientAddress, "- Topic:", topicName, "Parsed packet:", header.M_packetId, "->", string(jsonData))
	}
}

func parsePacketData(messageNoHeader []byte, packet any) (any, error) {
	err := binary.Read(bytes.NewReader(messageNoHeader), binary.LittleEndian, packet)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func main() {
	// Get server port from environment variable
	srv_port := os.Getenv("PORT")
	if srv_port == "" {
		log.Println("No PORT environment variable set, using default port", PORT)
		srv_port = PORT
	}

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

		go handleClientMessage(clientAddr, buffer[:n])
	}
}
