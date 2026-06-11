package mappers

import (
	mappersV23 "github.com/moon36/f1-game-telemetry/src/internal/mappers/v23"
	"github.com/moon36/f1-game-telemetry/src/internal/packets"
	packetsV23 "github.com/moon36/f1-game-telemetry/src/internal/packets/v23"
)

type PacketMapper interface {
	/*
		Maps a binary packet to the corresponding Go struct.

		Parameters:
		  - header: The packet header containing metadata about the packet.
		  - message: The binary data of the packet to be mapped.

		Returns:
		  - A pointer to a Go struct representing the mapped packet.
		  - A string for the associated message queue topic for this struct.
		  - An error if any occured.
	*/
	MapPacket(header packets.PacketHeader, message []byte) (any, string, error)
}

// Map to the correct mapper based on the packet format ID.
var PACKET_MAPPER_MAP = map[uint16]PacketMapper{
	packetsV23.PACKET_FORMAT_ID: &mappersV23.MapperV23{},
}
