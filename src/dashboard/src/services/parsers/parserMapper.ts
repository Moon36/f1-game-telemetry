import { Parser23 } from '@/services/parsers/telemetryParser23'
import type { TelemetryParser } from '@/types/telemetryParser'
import { PACKET_FORMAT_ID_23 } from '@/types/packets/v23/packetDefinitions'

export const parserMap: Record<number, TelemetryParser> = {
  [PACKET_FORMAT_ID_23]: new Parser23(),
}
