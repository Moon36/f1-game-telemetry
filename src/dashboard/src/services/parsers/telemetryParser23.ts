import type { TelemetryParser } from '@/types/telemetryParser';
import type { CarData } from '@/types/telemetry/carTelemetry';
import type { CarStatusData } from '@/types/telemetry/carStatusTelemetry';


export class Parser23 implements TelemetryParser {
  parseCarTelemetry(rawData: any, carIndex: number): CarData {
    const carData = rawData['M_carTelemetry']?.[carIndex] || {};
    return {
      innerTyreTemps: carData['M_tyresInnerTemperature'] || [0, 0, 0, 0],
      outerTyreTemps: carData['M_tyresSurfaceTemperature'] || [0, 0, 0, 0],
    };
  }

  parseCarStatusTelemetry(rawData: any, carIndex: number): CarStatusData {
    const statusData = rawData['M_carStatusData']?.[carIndex] || {};
    return {
      actualTyreCompoundId: statusData['M_actualTyreCompound'] ?? 0,
    };
  }
}