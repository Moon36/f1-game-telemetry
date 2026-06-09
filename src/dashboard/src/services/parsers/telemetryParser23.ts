import type { TelemetryParser } from '@/types/telemetryParser'
import type { CarData } from '@/types/telemetry/carTelemetry'
import type { CarStatusData } from '@/types/telemetry/carStatusTelemetry'
import type { RawTelemetry23 } from '@/types/packets/v23/packetDefinitions'
import type { M_carTelemetry, M_carStatusData } from '@/types/packets/v23/packetDefinitions'

export class Parser23 implements TelemetryParser<RawTelemetry23> {
  parseCarTelemetry(rawData: RawTelemetry23, carIndex: number): CarData {
    const carData = rawData['M_carTelemetry']?.[carIndex] as M_carTelemetry | undefined
    return {
      innerTyreTemps: carData?.M_tyresInnerTemperature || [0, 0, 0, 0],
      outerTyreTemps: carData?.M_tyresSurfaceTemperature || [0, 0, 0, 0],
    }
  }

  parseCarStatusTelemetry(rawData: RawTelemetry23, carIndex: number): CarStatusData {
    const statusData = rawData['M_carStatusData']?.[carIndex] as M_carStatusData | undefined
    return {
      actualTyreCompoundId: statusData?.M_actualTyreCompound ?? 0,
    }
  }
}
