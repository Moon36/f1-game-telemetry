import { Parser23 } from '@/services/parsers/telemetryParser23';
import type { TelemetryParser } from '@/types/telemetryParser';
import { PACKET_FORMAT_ID } from '@/constants';

export const parserMap: Record<number, TelemetryParser> = {
  [PACKET_FORMAT_ID]: new Parser23(),
};