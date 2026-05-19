import type { CarData } from '@/types/telemetry/carTelemetry';
import type { CarStatusData } from '@/types/telemetry/carStatusTelemetry';

export interface TelemetryParser {
  parseCarTelemetry(rawData: any, carIndex: number): CarData;
  parseCarStatusTelemetry(rawData: any, carIndex: number): CarStatusData;
}