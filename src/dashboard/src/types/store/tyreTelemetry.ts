export interface TyreTelemetry {
    innerTyreTemps: number[]; // [RL, RR, FL, FR]
    outerTyreTemps: number[]; // [RL, RR, FL, FR]
    actualTyreCompoundId: number; // 0 for unknown, otherwise 7-21
}