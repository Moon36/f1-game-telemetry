/*
 * Contains all packet information for F1 2023
 */
package packets

// === Packet ID to packet type map ===
var PACKET_MAP = map[uint8]any{
	CAR_MOTION_DATA_ID:           PacketMotionData{},
	SESSION_DATA_ID:              PacketSessionData{},
	LAP_DATA_ID:                  PacketLapData{},
	EVENT_DATA_ID:                GenericEvent{},
	PARTICIPANTS_DATA_ID:         PacketParticipantsData{},
	CAR_SETUP_DATA_ID:            PacketCarSetupData{},
	CAR_TELEMETRY_DATA_ID:        PacketCarTelemetryData{},
	CAR_STATUS_DATA_ID:           PacketCarStatusData{},
	FINAL_CLASSIFICATION_DATA_ID: PacketFinalClassificationData{},
	LOBBY_INFO_DATA_ID:           PacketLobbyInfoData{},
	CAR_DAMAGE_DATA_ID:           PacketCarDamageData{},
	SESSION_HISTORY_DATA_ID:      PacketSessionHistoryData{},
	TYRE_SET_DATA_ID:             PacketTyreSetsData{},
	CAR_MOTION_EX_DATA_ID:        PacketMotionExData{},
}
var PACKET_TOPIC_MAP = map[uint8]string{
	CAR_MOTION_DATA_ID:           TOPIC_CAR_MOTION_DATA,
	SESSION_DATA_ID:              TOPIC_SESSION_DATA,
	LAP_DATA_ID:                  TOPIC_LAP_DATA,
	EVENT_DATA_ID:                TOPIC_EVENT_DATA,
	PARTICIPANTS_DATA_ID:         TOPIC_PARTICIPANT_DATA,
	CAR_SETUP_DATA_ID:            TOPIC_CAR_SETUP_DATA,
	CAR_TELEMETRY_DATA_ID:        TOPIC_CAR_TELEMETRY_DATA,
	CAR_STATUS_DATA_ID:           TOPIC_CAR_STATUS_DATA,
	FINAL_CLASSIFICATION_DATA_ID: TOPIC_FINAL_CLASSIFICATION_DATA,
	LOBBY_INFO_DATA_ID:           TOPIC_LOBBY_INFO_DATA,
	CAR_DAMAGE_DATA_ID:           TOPIC_CAR_DAMAGE_DATA,
	SESSION_HISTORY_DATA_ID:      TOPIC_SESSION_HISTORY_DATA,
	TYRE_SET_DATA_ID:             TOPIC_TYRE_SET_DATA,
	CAR_MOTION_EX_DATA_ID:        TOPIC_CAR_MOTION_EX_DATA,
}

var EVENT_MAP = map[string]any{
	EVENT_CODE_SSTA: PacketEventSSTA{},
	EVENT_CODE_SEND: PacketEventSEND{},
	EVENT_CODE_FTLP: PacketEventFTLP{},
	EVENT_CODE_RTMT: PacketEventRTMT{},
	EVENT_CODE_DRSE: PacketEventDRSE{},
	EVENT_CODE_DRSD: PacketEventDRSD{},
	EVENT_CODE_TMPT: PacketEventTMPT{},
	EVENT_CODE_CHQF: PacketEventCHQF{},
	EVENT_CODE_RCWN: PacketEventRCWN{},
	EVENT_CODE_PENA: PacketEventPENA{},
	EVENT_CODE_SPTP: PacketEventSPTP{},
	EVENT_CODE_STLG: PacketEventSTLG{},
	EVENT_CODE_DTSV: PacketEventDTSV{},
	EVENT_CODE_SGSV: PacketEventSGSV{},
	EVENT_CODE_FLBK: PacketEventFLBK{},
	EVENT_CODE_BUTN: PacketEventBUTN{},
	EVENT_CODE_RDFL: PacketEventRDFL{},
	EVENT_CODE_OVTK: PacketEventOVTK{},
}

// === Packet IDs from the packetId field in the header ===
const CAR_MOTION_DATA_ID = 0
const SESSION_DATA_ID = 1
const LAP_DATA_ID = 2
const EVENT_DATA_ID = 3
const PARTICIPANTS_DATA_ID = 4
const CAR_SETUP_DATA_ID = 5
const CAR_TELEMETRY_DATA_ID = 6
const CAR_STATUS_DATA_ID = 7
const FINAL_CLASSIFICATION_DATA_ID = 8
const LOBBY_INFO_DATA_ID = 9
const CAR_DAMAGE_DATA_ID = 10
const SESSION_HISTORY_DATA_ID = 11
const TYRE_SET_DATA_ID = 12
const CAR_MOTION_EX_DATA_ID = 13

// === Packet sizes ===
const HEADER_SIZE = 24 // TODO Check if this is correct
const CAR_MOTION_DATA_SIZE = 1349
const SESSION_DATA_SIZE = 644
const LAP_DATA_SIZE = 1131
const EVENT_DATA_SIZE = 45
const PARTICIPANT_DATA_SIZE = 1306
const CAR_SETUP_DATA_SIZE = 1107
const CAR_TELEMETRY_DATA_SIZE = 1352
const CAR_STATUS_DATA_SIZE = 1239
const FINAL_CLASSIFICATION_DATA_SIZE = 1020
const LOBBY_INFO_DATA_SIZE = 1218
const CAR_DAMAGE_DATA_SIZE = 953
const SESSION_HISTORY_DATA_SIZE = 1460
const TYRE_SET_DATA_SIZE = 231
const CAR_MOTION_EXT_DATA_SIZE = 217

var MAX_BUFFER_SIZE = HEADER_SIZE + max(CAR_MOTION_DATA_SIZE, SESSION_DATA_SIZE, LAP_DATA_SIZE, EVENT_DATA_SIZE,
	PARTICIPANT_DATA_SIZE, CAR_SETUP_DATA_SIZE, CAR_TELEMETRY_DATA_SIZE, CAR_STATUS_DATA_SIZE,
	FINAL_CLASSIFICATION_DATA_SIZE, LOBBY_INFO_DATA_SIZE, CAR_DAMAGE_DATA_SIZE, SESSION_HISTORY_DATA_SIZE,
	TYRE_SET_DATA_SIZE, CAR_MOTION_EXT_DATA_SIZE)

// === Packet topics for Apache Kafka ===
const TOPIC_CAR_MOTION_DATA string = "telemetry.car_motion"
const TOPIC_SESSION_DATA string = "telemetry.session"
const TOPIC_LAP_DATA string = "telemetry.lap"
const TOPIC_EVENT_DATA string = "telemetry.event"
const TOPIC_PARTICIPANT_DATA string = "telemetry.participants"
const TOPIC_CAR_SETUP_DATA string = "telemetry.car_setup"
const TOPIC_CAR_TELEMETRY_DATA string = "telemetry.car_telemetry"
const TOPIC_CAR_STATUS_DATA string = "telemetry.car_status"
const TOPIC_FINAL_CLASSIFICATION_DATA string = "telemetry.final_classification"
const TOPIC_LOBBY_INFO_DATA string = "telemetry.lobby_info"
const TOPIC_CAR_DAMAGE_DATA string = "telemetry.car_damage"
const TOPIC_SESSION_HISTORY_DATA string = "telemetry.session_history"
const TOPIC_TYRE_SET_DATA string = "telemetry.tyre_set"
const TOPIC_CAR_MOTION_EX_DATA string = "telemetry.car_motion_ext"

var MESSAGE_TOPICS = [...]string{
	TOPIC_CAR_MOTION_DATA,
	TOPIC_SESSION_DATA,
	TOPIC_LAP_DATA,
	TOPIC_EVENT_DATA,
	TOPIC_PARTICIPANT_DATA,
	TOPIC_CAR_SETUP_DATA,
	TOPIC_CAR_TELEMETRY_DATA,
	TOPIC_CAR_STATUS_DATA,
	TOPIC_FINAL_CLASSIFICATION_DATA,
	TOPIC_LOBBY_INFO_DATA,
	TOPIC_CAR_DAMAGE_DATA,
	TOPIC_SESSION_HISTORY_DATA,
	TOPIC_TYRE_SET_DATA,
	TOPIC_CAR_MOTION_EX_DATA,
}

// === Event Codes ===
const EVENT_CODE_SSTA = "SSTA" // Session started
const EVENT_CODE_SEND = "SEND" // Session ended
const EVENT_CODE_FTLP = "FTLP" // Fastest lap
const EVENT_CODE_RTMT = "RTMT" // Retirement
const EVENT_CODE_DRSE = "DRSE" // DRS enabled
const EVENT_CODE_DRSD = "DRSD" // DRS disabled
const EVENT_CODE_TMPT = "TMPT" // Team mate in pits
const EVENT_CODE_CHQF = "CHQF" // Chequered flag
const EVENT_CODE_RCWN = "RCWN" // Race winner
const EVENT_CODE_PENA = "PENA" // Penalty issued
const EVENT_CODE_SPTP = "SPTP" // Speed trap triggered
const EVENT_CODE_STLG = "STLG" // Start lights
const EVENT_CODE_DTSV = "DTSV" // Drive through penalty served
const EVENT_CODE_SGSV = "SGSV" // Stop go penalty served
const EVENT_CODE_FLBK = "FLBK" // Flashback
const EVENT_CODE_BUTN = "BUTN" // Buttons
const EVENT_CODE_RDFL = "RDFL" // Red Flag
const EVENT_CODE_OVTK = "OVTK" // Overtake

// === Packet structs ===
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

/* Car Motion Data Packet
 * The motion packet gives physics data for all the cars being driven.
 *
 * Frequency: Rate as specified in the menus
 */
type PacketMotionData struct {
	//	m_header      PacketHeade
	CarMotionData [22]CarMotionData // Data for all cars
}
type CarMotionData struct {
	M_worldPositionX     float32 // World space X position - metres
	M_worldPositionY     float32 // World space Y position
	M_worldPositionZ     float32 // World space Z position
	M_worldVelocityX     float32 // Velocity in world space X – metres/s
	M_worldVelocityY     float32 // Velocity in world space Y
	M_worldVelocityZ     float32 // Velocity in world space Z
	M_worldForwardDirX   int16   // World space forward X direction (normalised)
	M_worldForwardDirY   int16   // World space forward Y direction (normalised)
	M_worldForwardDirZ   int16   // World space forward Z direction (normalised)
	M_worldRightDirX     int16   // World space right X direction (normalised)
	M_worldRightDirY     int16   // World space right Y direction (normalised)
	M_worldRightDirZ     int16   // World space right Z direction (normalised)
	M_gForceLateral      float32 // Lateral G-Force component
	M_gForceLongitudinal float32 // Longitudinal G-Force component
	M_gForceVertical     float32 // Vertical G-Force component
	M_yaw                float32 // Yaw angle in radians
	M_pitch              float32 // Pitch angle in radians
	M_roll               float32 // Roll angle in radians
}

/* Session Packet
 * The session packet includes details about the current session in progress.
 *
 * Frequency: 2 per second
 */
type PacketSessionData struct {
	//	m_header                          PacketHeader
	M_weather uint8 // Weather - 0 = clear, 1 = light cloud, 2 = overcast,
	// 3 = light rain, 4 = heavy rain, 5 = storm
	M_trackTemperature int8   // Track temp. in degrees celsius
	M_airTemperature   int8   // Air temp. in degrees celsius
	M_totalLaps        uint8  // Total number of laps in this race
	M_trackLength      uint16 // Track length in metres
	M_sessionType      uint8  // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P,
	// 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = R3, 13 = Time Trial
	M_trackId int8  // -1 for unknown, see appendix
	M_formula uint8 // Formula, 0 = F1 Modern, 1 = F1 Classic, 2 = F2,
	// 3 = F1 Generic, 4 = Beta, 5 = Supercars, 6 = Esports, 7 = F2 2021
	M_sessionTimeLeft     uint16          // Time left in session in seconds
	M_sessionDuration     uint16          // Session duration in seconds
	M_pitSpeedLimit       uint8           // Pit speed limit in kilometres per hour
	M_gamePaused          uint8           // Whether the game is paused – network game only
	M_isSpectating        uint8           // Whether the player is spectating
	M_spectatorCarIndex   uint8           // Index of the car being spectated
	M_sliProNativeSupport uint8           // SLI Pro support, 0 = inactive, 1 = active
	M_numMarshalZones     uint8           // Number of marshal zones to follow
	M_marshalZones        [21]MarshalZone // List of marshal zones – max 21
	M_safetyCarStatus     uint8           // 0 = no safety car, 1 = full, 2 = virtual,
	//  = formation lap
	M_networkGame               uint8                     // 0 = offline, 1 = online
	M_numWeatherForecastSamples uint8                     // Number of weather samples to follow
	M_weatherForecastSamples    [56]WeatherForecastSample // Array of weather forecast samples
	M_forecastAccuracy          uint8                     // 0 = Perfect, 1 = Approximate
	M_aiDifficulty              uint8                     // AI Difficulty rating – 0-110
	M_seasonLinkIdentifier      uint32                    // Identifier for season - persists across saves
	M_weekendLinkIdentifier     uint32                    // Identifier for weekend - persists across saves
	M_sessionLinkIdentifier     uint32                    // Identifier for session - persists across saves
	M_pitStopWindowIdealLap     uint8                     // Ideal lap to pit on for current strategy (player)
	M_pitStopWindowLatestLap    uint8                     // Latest lap to pit on for current strategy (player)
	M_pitStopRejoinPosition     uint8                     // Predicted position to rejoin at (player)
	M_steeringAssist            uint8                     // 0 = off, 1 = on
	M_brakingAssist             uint8                     // 0 = off, 1 = low, 2 = medium, 3 = high
	M_gearboxAssist             uint8                     // 1 = manual, 2 = manual & suggested gear, 3 = auto
	M_pitAssist                 uint8                     // 0 = off, 1 = on
	M_pitReleaseAssist          uint8                     // 0 = off, 1 = on
	M_ERSAssist                 uint8                     // 0 = off, 1 = on
	M_DRSAssist                 uint8                     // 0 = off, 1 = on
	M_dynamicRacingLine         uint8                     // 0 = off, 1 = corners only, 2 = full
	M_dynamicRacingLineType     uint8                     // 0 = 2D, 1 = 3D
	M_gameMode                  uint8                     // Game mode id - see appendix
	M_ruleSet                   uint8                     // Ruleset - see appendix
	M_timeOfDay                 uint32                    // Local time of day - minutes since midnight
	M_sessionLength             uint8                     // 0 = None, 2 = Very Short, 3 = Short, 4 = Medium,
	// 5 = Medium Long, 6 = Long, 7 = Full
	M_speedUnitsLeadPlayer            uint8 // 0 = MPH, 1 = KPH
	M_temperatureUnitsLeadPlayer      uint8 // 0 = Celsius, 1 = Fahrenheit
	M_speedUnitsSecondaryPlayer       uint8 // 0 = MPH, 1 = KPH
	M_temperatureUnitsSecondaryPlayer uint8 // 0 = Celsius, 1 = Fahrenheit
	M_numSafetyCarPeriods             uint8 // Number of safety cars called during session
	M_numVirtualSafetyCarPeriods      uint8 // Number of virtual safety cars called
	M_numRedFlagPeriods               uint8 // Number of red flags called during session
}
type MarshalZone struct {
	M_zoneStart float32 // Fraction (0..1) of way through the lap the marshal zone starts
	M_zoneFlag  int8    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow
}
type WeatherForecastSample struct {
	M_sessionType uint8 // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3,
	// 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = R3, 13 = Time Trial
	M_timeOffset uint8 // Time in minutes the forecast is for
	M_weather    uint8 // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain,
	// 4 = heavy rain, 5 = storm
	M_trackTemperature       int8  // Track temp. in degrees Celsius
	M_trackTemperatureChange int8  // Track temp. change – 0 = up, 1 = down, 2 = no change
	M_airTemperature         int8  // Air temp. in degrees celsius
	M_airTemperatureChange   int8  // Air temp. change – 0 = up, 1 = down, 2 = no change
	M_rainPercentage         uint8 // Rain percentage (0-100)
}

/* Lap Packet
 * The lap data packet gives details of all the cars in the session.
 *
 * Frequency: Rate as specified in the menus
 */
type PacketLapData struct {
	//	m_header               PacketHeader
	M_lapData              [22]LapData // Lap data for all cars on track
	M_timeTrialPBCarIdx    uint8       // Index of Personal Best car in time trial (255 if invalid)
	M_timeTrialRivalCarIdx uint8       // Index of Rival car in time trial (255 if invalid)
}
type LapData struct {
	M_lastLapTimeInMS       uint32  // Last lap time in milliseconds
	M_currentLapTimeInMS    uint32  // Current time around the lap in milliseconds
	M_sector1TimeInMS       uint16  // Sector 1 time in milliseconds
	M_sector1TimeMinutes    uint8   // Sector 1 whole minute part
	M_sector2TimeInMS       uint16  // Sector 2 time in milliseconds
	M_sector2TimeMinutes    uint8   // Sector 2 whole minute part
	M_deltaToCarInFrontInMS uint16  // Time delta to car in front in milliseconds
	M_deltaToRaceLeaderInMS uint16  // Time delta to race leader in milliseconds
	M_lapDistance           float32 // Distance vehicle is around current lap in metres – could be negative if
	// line hasn’t been crossed yet
	M_totalDistance float32 // Total distance travelled in session in metres – could be negative if line
	// hasn’t been crossed yet
	M_safetyCarDelta              float32 // Delta in seconds for safety car
	M_carPosition                 uint8   // Car race position
	M_currentLapNum               uint8   // Current lap number
	M_pitStatus                   uint8   // 0 = none, 1 = pitting, 2 = in pit area
	M_numPitStops                 uint8   // Number of pit stops taken in this race
	M_sector                      uint8   // 0 = sector1, 1 = sector2, 2 = sector3
	M_currentLapInvalid           uint8   // Current lap invalid - 0 = valid, 1 = invalid
	M_penalties                   uint8   // Accumulated time penalties in seconds to be added
	M_totalWarnings               uint8   // Accumulated number of warnings issued
	M_cornerCuttingWarnings       uint8   // Accumulated number of corner cutting warnings issued
	M_numUnservedDriveThroughPens uint8   // Num drive through pens left to serve
	M_numUnservedStopGoPens       uint8   // Num stop go pens left to serve
	M_gridPosition                uint8   // Grid position the vehicle started the race in
	M_driverStatus                uint8   // Status of driver - 0 = in garage, 1 = flying lap 2 = in lap, 3 = out lap,
	// 4 = on track
	M_resultStatus uint8 // Result status - 0 = invalid, 1 = inactive, 2 = active 3 = finished,
	// 4 = didnotfinish, 5 = disqualified 6 = not classified, 7 = retired
	M_pitLaneTimerActive    uint8  // Pit lane timing, 0 = inactive, 1 = active
	M_pitLaneTimeInLaneInMS uint16 // If active, the current time spent in the pit lane in ms
	M_pitStopTimerInMS      uint16 // Time of the actual pit stop in ms
	M_pitStopShouldServePen uint8  // Whether the car should serve a penalty at this stop
}

/* Event Packet
 * This packet gives details of events that happen during the course of a session.
 *
 * Frequency: When the event occurs
 */
type GenericEvent struct {
	M_eventStringCode [4]uint8
	M_eventDetails    any
}

// Session started
type PacketEventSSTA struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for SSTA
}

// Session ended
type PacketEventSEND struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for SEND
}

// Fastest lap
type PacketEventFTLP struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    FastestLap // Details of the fastest lap event
}

// Retirement
type PacketEventRTMT struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    Retirement // Details of the retirement event
}

// DRS enabled
type PacketEventDRSE struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for DRSE
}

// DRS disabled
type PacketEventDRSD struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for DRSD
}

// Team mate in pits
type PacketEventTMPT struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    TeamMateInPits // Details of the team mate in pits event
}

// Chequered flag
type PacketEventCHQF struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for CHQF
}

// Race winner
type PacketEventRCWN struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    RaceWinner // Details of the race winner event
}

// Penalty issued
type PacketEventPENA struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    Penalty // Details of the penalty event
}

// Speed trap triggered
type PacketEventSPTP struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    SpeedTrap // Details of the speed trap event
}

// Start lights
type PacketEventSTLG struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    StartLights // Details of the start lights event
}

// Drive through penalty served
type PacketEventDTSV struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    DriveThroughPenaltyServed // Details of the drive through penalty served event
}

// Stop go penalty served
type PacketEventSGSV struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    StopGoPenaltyServed // Details of the stop go penalty served event
}

// Flashback
type PacketEventFLBK struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    Flashback // Details of the flashback event
}

// Buttons
type PacketEventBUTN struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    Buttons // Details of the buttons event
}

// Red Flag
type PacketEventRDFL struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    any // No event details for RDFL
}

// Overtake
type PacketEventOVTK struct {
	//	m_header          PacketHeader
	M_eventStringCode [4]uint8
	M_eventDetails    Overtake // Details of the overtake event
}

// Define event type structs
type FastestLap struct {
	vehicleIdx uint8   // Vehicle index of car achieving fastest lap
	lapTime    float32 // Lap time is in seconds
}
type Retirement struct {
	vehicleIdx uint8 // Vehicle index of car retiring
}
type TeamMateInPits struct {
	vehicleIdx uint8 // Vehicle index of team mate
}
type RaceWinner struct {
	vehicleIdx uint8 // Vehicle index of the race winner
}
type Penalty struct {
	penaltyType      uint8 // Penalty type – see Appendices
	infringementType uint8 // Infringement type – see Appendices
	vehicleIdx       uint8 // Vehicle index of the car the penalty is applied to
	otherVehicleIdx  uint8 // Vehicle index of the other car involved
	time             uint8 // Time gained, or time spent doing action in seconds
	lapNum           uint8 // Lap the penalty occurred on
	placesGained     uint8 // Number of places gained by this
}
type SpeedTrap struct {
	vehicleIdx                 uint8   // Vehicle index of the vehicle triggering speed trap
	speed                      float32 // Top speed achieved in kilometres per hour
	isOverallFastestInSession  uint8   // Overall fastest speed in session = 1, otherwise 0
	isDriverFastestInSession   uint8   // Fastest speed for driver in session = 1, otherwise 0
	fastestVehicleIdxInSession uint8   // Vehicle index of the vehicle that is the fastest in this session
	fastestSpeedInSession      float32 // Speed of the vehicle that is the fastest in this session
}
type StartLights struct {
	numLights uint8 // Number of lights showing
}
type DriveThroughPenaltyServed struct {
	vehicleIdx uint8 // Vehicle index of the vehicle serving drive through
}
type StopGoPenaltyServed struct {
	vehicleIdx uint8 // Vehicle index of the vehicle serving stop go
}
type Flashback struct {
	flashbackFrameIdentifier uint32  // Frame identifier flashed back to
	flashbackSessionTime     float32 // Session time flashed back to
}
type Buttons struct {
	buttonStatus uint32 // Bit flags specifying which buttons are being pressed currently - see appendices
}
type Overtake struct {
	overtakingVehicleIdx     uint8 // Vehicle index of the vehicle overtaking
	beingOvertakenVehicleIdx uint8 // Vehicle index of the vehicle being overtaken
}

/* Participant Packet
 * This is a list of participants in the race. If the vehicle is controlled by AI, then the name will be the driver
 * name. If this is a multiplayer game, the names will be the Steam Id on PC, or the LAN name if appropriate.
 * N.B. on Xbox One, the names will always be the driver name, on PS4 the name will be the LAN name if playing a LAN
 * game, otherwise it will be the driver name.
 * The array should be indexed by vehicle index.
 *
 * Frequency: Every 5 seconds
 */
type PacketParticipantsData struct {
	//	m_header        PacketHeader
	M_numActiveCars uint8               // Number of active cars in the data
	M_participants  [22]ParticipantData // Data for all participants
}
type ParticipantData struct {
	M_aiControlled uint8    // Whether the vehicle is AI (1) or Human (0) controlled
	M_driverId     uint8    // Driver id - see appendix, 255 if network human
	M_networkId    uint8    // Network id – unique identifier for network players
	M_teamId       uint8    // Team id - see appendix
	M_myTeam       uint8    // My team flag – 1 = My Team, 0 = otherwise
	M_raceNumber   uint8    // Race number of the car
	M_nationality  uint8    // Nationality of the driver
	M_name         [48]byte // Name of participant in UTF-8 format – null terminated. Will be truncated with "..."
	// (U+2026) if too long
	M_yourTelemetry   uint8 // The player's UDP setting, 0 = restricted, 1 = public
	M_showOnlineNames uint8 // The player's show online names setting, 0 = off, 1 = on
	M_platform        uint8 // 1 = Steam, 3 = PlayStation, 4 = Xbox, 6 = Origin, 255 = unknown
}

/* Car Setups Packet
 * This packet details the car setups for each vehicle in the session. Note that in multiplayer games, other player cars
 * will appear as blank, you will only be able to see your own car setup, regardless of the “Your Telemetry” setting.
 * Spectators will also not be able to see any car setups.
 *
 * Frequency: 2 per second
 */
type PacketCarSetupData struct {
	//	m_header    PacketHeader
	M_carSetups [22]CarSetupData
}
type CarSetupData struct {
	M_frontWing              uint8   // Front wing aero
	M_rearWing               uint8   // Rear wing aero
	M_onThrottle             uint8   // Differential adjustment on throttle (percentage)
	M_offThrottle            uint8   // Differential adjustment off throttle (percentage)
	M_frontCamber            float32 // Front camber angle (suspension geometry)
	M_rearCamber             float32 // Rear camber angle (suspension geometry)
	M_frontToe               float32 // Front toe angle (suspension geometry)
	M_rearToe                float32 // Rear toe angle (suspension geometry)
	M_frontSuspension        uint8   // Front suspension
	M_rearSuspension         uint8   // Rear suspension
	M_frontAntiRollBar       uint8   // Front anti-roll bar
	M_rearAntiRollBar        uint8   // Front anti-roll bar
	M_frontSuspensionHeight  uint8   // Front ride height
	M_rearSuspensionHeight   uint8   // Rear ride height
	M_brakePressure          uint8   // Brake pressure (percentage)
	M_brakeBias              uint8   // Brake bias (percentage)
	M_rearLeftTyrePressure   float32 // Rear left tyre pressure (PSI)
	M_rearRightTyrePressure  float32 // Rear right tyre pressure (PSI)
	M_frontLeftTyrePressure  float32 // Front left tyre pressure (PSI)
	M_frontRightTyrePressure float32 // Front right tyre pressure (PSI)
	M_ballast                uint8   // Ballast
	M_fuelLoad               float32 // Fuel load
}

/* Car Telemetry Packet
 * This packet details telemetry for all the cars in the race. It details various values that would be recorded on the
 * car such as speed, throttle application, DRS etc. Note that the rev light configurations are presented separately as
 * well and will mimic real life driver preferences.
 */
type PacketCarTelemetryData struct {
	//	m_header                       PacketHeader
	M_carTelemetry  [22]CarTelemetryData
	M_mfdPanelIndex uint8 // Index of the MFD panel open, 255 if closed; Single player race:
	// 0 = Car setup, 1 = Pits, 2 = Damage, 3 = Engine, 4 = Temperatures
	M_mfdPanelIndexSecondaryPlayer uint8 // see above
	M_suggestedGear                uint8 // Suggested gear for the player (1 - 8), 0 if no gear suggested
}
type CarTelemetryData struct {
	M_speed                   uint16     // Speed of car in kilometres per hour
	M_throttle                float32    // Amount of throttle applied (0.0 to 1.0)
	M_steer                   float32    // Steering (-1.0 (full lock left) to 1.0 (full lock right))
	M_brake                   float32    // Amount of brake applied (0.0 to 1.0)
	M_clutch                  uint8      // Amount of clutch applied (0 to 100)
	M_gear                    int8       // Gear selected (1-8, N=0, R=-1)
	M_engineRPM               uint16     // Engine RPM
	M_drs                     uint8      // 0 = off, 1 = on
	M_revLightsPercent        uint8      // Rev lights indicator (percentage)
	M_revLightsBitValue       uint16     // Rev lights (bit 0 = leftmost LED, bit 14 = rightmost LED)
	M_brakesTemperature       [4]uint16  // Brakes temperature (celsius)
	M_tyresSurfaceTemperature [4]uint8   // Tyres surface temperature (celsius)
	M_tyresInnerTemperature   [4]uint8   // Tyres inner temperature (celsius)
	M_engineTemperature       uint16     // Engine temperature (celsius)
	M_tyresPressure           [4]float32 // Tyres pressure (PSI)
	M_surfaceType             [4]uint8   // Driving surface
}

/* Car Status Packet
 * This packet details car statuses for all the cars in the race.
 *
 * Frequency: Rate as specified in the menus
 */
type PacketCarStatusData struct {
	//	m_header        PacketHeader
	M_carStatusData [22]CarStatusData // Car status data for all cars
}
type CarStatusData struct {
	M_tractionControl       uint8   // Traction control - 0 = off, 1 = medium, 2 = full
	M_antiLockBrakes        uint8   // 0 (off) - 1 (on)
	M_fuelMix               uint8   // Fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
	M_frontBrakeBias        uint8   // Front brake bias (percentage)
	M_pitLimiterStatus      uint8   // Pit limiter status - 0 = off, 1 = on
	M_fuelInTank            float32 // Current fuel mass
	M_fuelCapacity          float32 // Fuel capacity
	M_fuelRemainingLaps     float32 // Fuel remaining in terms of laps (value on MFD)
	M_maxRPM                uint16  // Cars max RPM, point of rev limiter
	M_idleRPM               uint16  // Cars idle RPM
	M_maxGears              uint8   // Maximum number of gears
	M_drsAllowed            uint8   // 0 = not allowed, 1 = allowed
	M_drsActivationDistance uint16  // 0 = DRS not available, non-zero - DRS will be available, in [X] metres
	M_actualTyreCompound    uint8   // F1 Modern - 16 = C5, 17 = C4, 18 = C3, 19 = C2, 20 = C1, 21 = C0, 7 = inter,
	// 8 = wet, F1 Classic - 9 = dry, 10 = wet, F2 – 11 = super soft, 12 = soft, 13 = medium, 14 = hard, 15 = wet
	M_visualTyreCompound uint8 // F1 visual (can be different from actual compound), 16 = soft, 17 = medium,
	// 18 = hard, 7 = inter, 8 = wet, F1 Classic – same as above, F2 ‘19, 15 = wet, 19 – super soft, 20 = soft,
	// 21 = medium , 22 = hard
	M_tyresAgeLaps            uint8   // Age in laps of the current set of tyres
	M_vehicleFiaFlags         int8    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow
	M_enginePowerICE          float32 // Engine power output of ICE (W)
	M_enginePowerMGUK         float32 // Engine power output of MGU-K (W)
	M_ersStoreEnergy          float32 // ERS energy store in Joules
	M_ersDeployMode           uint8   // ERS deployment mode, 0 = none, 1 = medium, 2 = hotlap, 3 = overtake
	M_ersHarvestedThisLapMGUK float32 // ERS energy harvested this lap by MGU-K
	M_ersHarvestedThisLapMGUH float32 // ERS energy harvested this lap by MGU-H
	M_ersDeployedThisLap      float32 // ERS energy deployed this lap
	M_networkPaused           uint8   // Whether the car is paused in a network game
}

/* Final Classification Packet
 * This packet details the final classification at the end of the race, and the data will match with the post race
 * results screen. This is especially useful for multiplayer games where it is not always possible to send lap times on
 * the final frame because of network delay.
 *
 * Frequency: Once at the end of the race
 */
type PacketFinalClassificationData struct {
	//	m_header             PacketHeader
	M_numCars            uint8                       // Number of cars in the final classification
	M_classificationData [22]FinalClassificationData // Final classification data for all cars
}
type FinalClassificationData struct {
	M_position     uint8 // Finishing position
	M_numLaps      uint8 // Number of laps completed
	M_gridPosition uint8 // Grid position of the car
	M_points       uint8 // Number of points scored
	M_numPitStops  uint8 // Number of pit stops made
	M_resultStatus uint8 // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished,
	// 4 = didnotfinish, 5 = disqualified, 6 = not classified, 7 = retired
	M_bestLapTimeInMS   uint32   // Best lap time of the session in milliseconds
	M_totalRaceTime     float64  // Total race time in seconds without penalties
	M_penaltiesTime     uint8    // Total penalties accumulated in seconds
	M_numPenalties      uint8    // Number of penalties applied to this driver
	M_numTyreStints     uint8    // Number of tyres stints up to maximum
	M_tyreStintsActual  [8]uint8 // Actual tyres used by this driver
	M_tyreStintsVisual  [8]uint8 // Visual tyres used by this driver
	M_tyreStintsEndLaps [8]uint8 // The lap number stints end on
}

/* Lobby Info Packet
 * This packet details the players currently in a multiplayer lobby. It details each player’s selected car, any AI
 * involved in the game and also the ready status of each of the participants.
 *
 * Frequency: 2 per second when in a lobby
 */
type PacketLobbyInfoData struct {
	//	m_header       PacketHeader
	M_numPlayers   uint8 // Number of players in the lobby
	M_lobbyPlayers [22]LobbyPlayerData
}
type LobbyPlayerData struct {
	M_aiControlled uint8    // Whether the vehicle is AI (1) or Human (0) controlled
	M_teamId       uint8    // Team id - see appendix (255 if no team currently selected)
	M_nationality  uint8    // Nationality of the driver
	M_platform     uint8    // 1 = Steam, 3 = PlayStation, 4 = Xbox, 6 = Origin, 255 = unknown
	M_name         [48]byte // Name of participant in UTF-8 format – null terminated will be truncated with ... (U+2026)
	//  if too long
	M_carNumber   uint8 // Car number of the player
	M_readyStatus uint8 // 0 = not ready, 1 = ready, 2 = spectating
}

/* Car Damage Packet
 * This packet details car damage parameters for all the cars in the race.
 *
 * Frequency: 10 per second
 */
type PacketCarDamageData struct {
	//	m_header        PacketHeader
	M_carDamageData [22]CarDamageData
}
type CarDamageData struct {
	M_tyresWear            [4]float32 // Tyre wear (percentage)
	M_tyresDamage          [4]uint8   // Tyre damage (percentage)
	M_brakesDamage         [4]uint8   // Brakes damage (percentage)
	M_frontLeftWingDamage  uint8      // Front left wing damage (percentage)
	M_frontRightWingDamage uint8      // Front right wing damage (percentage)
	M_rearWingDamage       uint8      // Rear wing damage (percentage)
	M_floorDamage          uint8      // Floor damage (percentage)
	M_diffuserDamage       uint8      // Diffuser damage (percentage)
	M_sidepodDamage        uint8      // Sidepod damage (percentage)
	M_drsFault             uint8      // Indicator for DRS fault, 0 = OK, 1 = fault
	M_ersFault             uint8      // Indicator for ERS fault, 0 = OK, 1 = fault
	M_gearBoxDamage        uint8      // Gear box damage (percentage)
	M_engineDamage         uint8      // Engine damage (percentage)
	M_engineMGUHWear       uint8      // Engine wear MGU-H (percentage)
	M_engineESWear         uint8      // Engine wear ES (percentage)
	M_engineCEWear         uint8      // Engine wear CE (percentage)
	M_engineICEWear        uint8      // Engine wear ICE (percentage)
	M_engineMGUKWear       uint8      // Engine wear MGU-K (percentage)
	M_engineTCWear         uint8      // Engine wear TC (percentage)
	M_engineBlown          uint8      // Engine blown, 0 = OK, 1 = fault
	M_engineSeized         uint8      // Engine seized, 0 = OK, 1 = fault
}

/* Session History Packet
 * This packet contains lap times and tyre usage for the session. This packet works slightly differently to other
 * packets. To reduce CPU and bandwidth, each packet relates to a specific vehicle and is sent every 1/20 s, and the
 * vehicle being sent is cycled through. Therefore in a 20 car race you should receive an update for each vehicle at
 * least once per second.
 * Note that at the end of the race, after the final classification packet has been sent, a final bulk update of all the
 * session histories for the vehicles in that session will be sent.
 *
 * Frequency: 20 per second but cycling through cars
 */
type PacketSessionHistoryData struct {
	//	m_header                PacketHeader
	M_carIdx                uint8               // Index of the car this lap data relates to
	M_numLaps               uint8               // Num laps in the data (including current partial lap)
	M_numTyreStints         uint8               // Number of tyre stints in the data
	M_bestLapTimeLapNum     uint8               // Lap the best lap time was achieved on
	M_bestSector1LapNum     uint8               // Lap the best Sector 1 time was achieved on
	M_bestSector2LapNum     uint8               // Lap the best Sector 2 time was achieved on
	M_bestSector3LapNum     uint8               // Lap the best Sector 3 time was achieved on
	M_lapHistoryData        [100]LapHistoryData // 100 laps of data max
	M_tyreStintsHistoryData [8]TyreStintHistoryData
}
type LapHistoryData struct {
	M_lapTimeInMS        uint32 // Lap time in milliseconds
	M_sector1TimeInMS    uint16 // Sector 1 time in milliseconds
	M_sector1TimeMinutes uint8  // Sector 1 whole minute part
	M_sector2TimeInMS    uint16 // Sector 2 time in milliseconds
	M_sector2TimeMinutes uint8  // Sector 2 whole minute part
	M_sector3TimeInMS    uint16 // Sector 3 time in milliseconds
	M_sector3TimeMinutes uint8  // Sector 3 whole minute part
	M_lapValidBitFlags   uint8  // 0x01 bit set-lap valid, 0x02 bit set-sector 1 valid, 0x04 bit set-sector 2 valid,
	// 0x08 bit set-sector 3 valid
}

type TyreStintHistoryData struct {
	M_endLap             uint8 // Lap the tyre usage ends on (255 of current tyre)
	M_tyreActualCompound uint8 // Actual tyres used by this driver
	M_tyreVisualCompound uint8 // Visual tyres used by this driver
}

/* Tyre Sets Packet
 * This packets gives a more in-depth details about tyre sets assigned to a vehicle during the session.
 *
 * Frequency: 20 per second but cycling through cars
 */
type PacketTyreSetsData struct {
	//	m_header    PacketHeader
	M_carIdx    uint8           // Index of the car this data relates to
	M_tyreSets  [20]TyreSetData // 13 (dry) + 7 (wet)
	M_fittedIdx uint8           // Index into array of fitted tyre
}

type TyreSetData struct {
	M_actualTyreCompound uint8 // Actual tyre compound used
	M_visualTyreCompound uint8 // Visual tyre compound used
	M_wear               uint8 // Tyre wear (percentage)
	M_available          uint8 // Whether this set is currently available
	M_recommendedSession uint8 // Recommended session for tyre set
	M_lifeSpan           uint8 // Laps left in this tyre set
	M_usableLife         uint8 // Max number of laps recommended for this compound
	M_lapDeltaTime       int16 // Lap delta time in milliseconds compared to fitted set
	M_fitted             uint8 // Whether the set is fitted or not
}

/* Motion Ex Packet
 * The motion packet gives extended data for the car being driven with the goal of being able to drive a motion platform
 * setup.
 *
 * Frequency: Rate as specified in the menus
 */
type PacketMotionExData struct {
	//	m_header PacketHeader // Header
	// Extra player car ONLY data
	M_suspensionPosition     [4]float32 // Note: All wheel arrays have the following order:
	M_suspensionVelocity     [4]float32 // RL, RR, FL, FR
	M_suspensionAcceleration [4]float32 // RL, RR, FL, FR
	M_wheelSpeed             [4]float32 // Speed of each wheel
	M_wheelSlipRatio         [4]float32 // Slip ratio for each wheel
	M_wheelSlipAngle         [4]float32 // Slip angles for each wheel
	M_wheelLatForce          [4]float32 // Lateral forces for each wheel
	M_wheelLongForce         [4]float32 // Longitudinal forces for each wheel
	M_heightOfCOGAboveGround float32    // Height of centre of gravity above ground
	M_localVelocityX         float32    // Velocity in local space – metres/s
	M_localVelocityY         float32    // Velocity in local space
	M_localVelocityZ         float32    // Velocity in local space
	M_angularVelocityX       float32    // Angular velocity x-component – radians/s
	M_angularVelocityY       float32    // Angular velocity y-component
	M_angularVelocityZ       float32    // Angular velocity z-component
	M_angularAccelerationX   float32    // Angular acceleration x-component – radians/s/s
	M_angularAccelerationY   float32    // Angular acceleration y-component
	M_angularAccelerationZ   float32    // Angular acceleration z-component
	M_frontWheelsAngle       float32    // Current front wheels angle in radians
	M_wheelVertForce         [4]float32 // Vertical forces for each wheel
}
