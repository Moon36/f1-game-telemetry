package v23

import (
	"encoding/binary"
	"fmt"

	"github.com/moon36/f1-game-telemetry/src/internal/packets"
	packetsV23 "github.com/moon36/f1-game-telemetry/src/internal/packets/v23"
)

type MapperV23 struct{}

func (m *MapperV23) MapPacket(header packets.PacketHeader, message []byte) (any, string, error) {
	packetId := header.M_packetId
	createPacket, prst := packetsV23.PACKET_MAP[packetId]
	if !prst {
		return nil, "", fmt.Errorf("unknown packet ID: %d", packetId)
	}

	headerSize := binary.Size(header)
	messageNoHeader := message[headerSize:]
	var packetData any
	var topicName string

	if packetId == packetsV23.EVENT_DATA_ID {
		if len(messageNoHeader) < 4 {
			return nil, "", fmt.Errorf("supposed event data packet ('23) does not contain event code")
		}
		eventCode := string(messageNoHeader[:4])
		createPacket, prst = packetsV23.EVENT_MAP[eventCode]
		if !prst {
			return nil, "", fmt.Errorf("unknown event code: %s", eventCode)
		}
	}
	packetData = createPacket()
	topicName = packetsV23.PACKET_TOPIC_MAP[packetId]

	return packetData, topicName, nil
}
