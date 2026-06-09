/*
 * Contains constants and types used across the database manager component.
 */
package main

import (
	"encoding/json"

	"github.com/moon36/f1-game-telemetry/src/internal/packets"
)

const USERNAME = "admin"
const PASSWORD = "MyPassword"

const KAFKA_CONSUMER_GROUP_ID = "db-manager-consumer"

// BasePacket to parse the header of incoming packets
type BasePacket struct {
	M_header packets.PacketHeader
	M_data   json.RawMessage
}

// ChannelData is the struct used to send enriched participant data along with any potential error through a channel
type ChannelData struct {
	Idx             uint8
	ParticipantData packets.StoreParticipant
	Error           error
}
