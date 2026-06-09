import type { GenericRawTelemetry } from '@/types/packets/packetDefinitions'

export const PACKET_FORMAT_ID_23 = 2023

export interface RawTelemetry23 extends GenericRawTelemetry {
  M_carTelemetry: M_carTelemetry[]
  M_carStatusData: M_carStatusData[]
}

export interface M_carTelemetry {
  M_speed: number
  M_throttle: number
  M_steer: number
  M_brake: number
  M_clutch: number
  M_gear: number
  M_engineRPM: number
  M_drs: number
  M_revLightsPercent: number
  M_revLightsBitValue: number
  M_brakesTemperature: number[]
  M_tyresSurfaceTemperature: number[]
  M_tyresInnerTemperature: number[]
  M_engineTemperature: number
  M_tyresPressure: number[]
  M_surfaceType: number[]
}

export interface M_carStatusData {
  M_tractionControl: number
  M_antiLockBrakes: number
  M_fuelMix: number
  M_frontBrakeBias: number
  M_pitLimiterStatus: number
  M_fuelInTank: number
  M_fuelCapacity: number
  M_fuelRemainingLaps: number
  M_maxRPM: number
  M_idleRPM: number
  M_maxGears: number
  M_drsAllowed: number
  M_drsActivationDistance: number
  M_actualTyreCompound: number
  M_visualTyreCompound: number
  M_tyresAgeLaps: number
  M_vehicleFiaFlags: number
  M_enginePowerICE: number
  M_enginePowerMGUK: number
  M_ersStoreEnergy: number
  M_ersDeployMode: number
  M_ersHarvestedThisLapMGUK: number
  M_ersHarvestedThisLapMGUH: number
  M_ersDeployedThisLap: number
  M_networkPaused: number
}
