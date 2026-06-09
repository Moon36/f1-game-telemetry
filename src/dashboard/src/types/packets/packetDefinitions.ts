// Generic, version-agnostic raw telemetry structure.
export type GenericRawTelemetry<
  PacketKeys extends string = string,
  PacketEntries extends object = object,
> = Record<PacketKeys, PacketEntries[]>
