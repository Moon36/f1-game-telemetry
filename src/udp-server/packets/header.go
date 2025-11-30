/*
 * Contains the definition of the packet header structure used in the F1 game telemetry data.
 */
package packets

/*
 * Each packet carries different types of data rather than having one packet which contains everything. The header in
 * each packet describes the packet type and versioning info so it will be easier for applications to check they are
 * interpreting the incoming data in the correct way. Please note that all values are encoded using Little Endian
 * format. All data is packed.
 */
type PacketHeader struct {
	M_packetFormat           uint16  // 2023
	M_gameYear               uint8   // Game year - last two digits e.g. 23
	M_gameMajorVersion       uint8   // Game major version - "X.00"
	M_gameMinorVersion       uint8   // Game minor version - "1.XX"
	M_packetVersion          uint8   // Version of this packet type, all start from 1
	M_packetId               uint8   // Identifier for the packet type, see below
	M_sessionUID             uint64  // Unique identifier for the session
	M_sessionTime            float32 // Session timestamp
	M_frameIdentifier        uint32  // Identifier for the frame the data was retrieved on
	M_overallFrameIdentifier uint32  // Overall identifier for the frame the data was retrieved on, doesn't go back
	// after flashbacks
	M_playerCarIndex          uint8 // Index of player's car in the array
	M_secondaryPlayerCarIndex uint8 // Index of secondary player's car in the array (splitscreen) 255 if no second
	// player
}
