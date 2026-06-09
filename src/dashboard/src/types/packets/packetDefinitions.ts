// Generic, version-agnostic raw telemetry structure.
// packetDefinitions.ts

export interface GenericRawTelemetry {
  M_header: Header
  [key: string]: unknown
}

export interface Header {
  M_packetFormat: number // e.g. 2023
  M_gameYear: number // Game year - last two digits e.g. 23
  M_gameMajorVersion: number // Game major version - "X.00"
  M_gameMinorVersion: number // Game minor version - "1.XX"
  M_packetVersion: number // Version of this packet type, all start from 1
  M_packetId: number // Identifier for the packet type, see below
  M_sessionUID: number // Unique identifier for the session
  M_sessionTime: number // Session timestamp
  M_frameIdentifier: number // Identifier for the frame the data was retrieved on
  M_overallFrameIdentifier: number // Overall identifier for the frame the data was retrieved on, doesn't go back after
  // flashbacks
  M_playerCarIndex: number // Index of player's car in the array
  M_secondaryPlayerCarIndex: number // Index of secondary player's car in the array (splitscreen) 255 if no second
  // player
}
