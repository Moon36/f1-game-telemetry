/*
 * Contains all common constants across all versions of the F1 game telemetry packets.
 */

package packets

const MAX_BUFFER_SIZE = 2048 // Maximum buffer size for incoming UDP packets, manually set

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
