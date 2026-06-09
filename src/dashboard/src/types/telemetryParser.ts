import type { GenericRawTelemetry } from '@/types/packets/packetDefinitions'
import type { CarData } from '@/types/telemetry/carTelemetry'
import type { CarStatusData } from '@/types/telemetry/carStatusTelemetry'

export interface TelemetryParser<RawTelemetry = GenericRawTelemetry> {
  parseCarTelemetry(rawData: RawTelemetry, carIndex: number): CarData
  parseCarStatusTelemetry(rawData: RawTelemetry, carIndex: number): CarStatusData
}
