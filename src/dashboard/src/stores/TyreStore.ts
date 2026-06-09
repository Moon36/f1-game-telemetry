import type { TyreTelemetry } from '@/types/store/tyreTelemetry'
import { TyreCompounds } from '@/utils/constants'
import { defineStore } from 'pinia'

export const useTyreStore = defineStore('tyreStore', {
  state: () =>
    ({
      innerTyreTemps: [0, 0, 0, 0],
      outerTyreTemps: [0, 0, 0, 0],
      actualTyreCompoundId: 0,
    }) as TyreTelemetry,
  actions: {
    /**
     * Updates the inner and outer tyre temperatures.
     *
     * @param innerTemps - Array of four numbers representing the inner tyre temperatures.
     * @param outerTemps - Array of four numbers representing the outer tyre temperatures.
     */
    updateTyreTemps(innerTemps: number[], outerTemps: number[]) {
      if (innerTemps.length === 4) {
        this.innerTyreTemps.splice(0, 4, ...innerTemps)
      } else {
        console.warn('Invalid inner tyre temperature data:', innerTemps)
      }
      if (outerTemps.length === 4) {
        this.outerTyreTemps.splice(0, 4, ...outerTemps)
      } else {
        console.warn('Invalid outer tyre temperature data:', outerTemps)
      }
    },
    /**
     * Updates the actual tyre compound ID.
     *
     * @param compoundId - Number representing the tyre compound ID.
     */
    updateTyreCompound(compoundId: number) {
      if (typeof compoundId === 'number') {
        this.actualTyreCompoundId = compoundId
      } else {
        console.warn('Invalid tyre compound ID data:', compoundId)
      }
    },
  },
  getters: {
    /**
     * Maps the compound ID from telemetry data to TyreCompounds enum.
     *
     * @returns The corresponding TyreCompounds enum value.
     */
    tyreCompound(): TyreCompounds {
      switch (this.actualTyreCompoundId) {
        case 7:
          return TyreCompounds.INTERMEDIATE
        case 8:
          return TyreCompounds.WET
        case 9:
          return TyreCompounds.DRY
        case 10:
          return TyreCompounds.WET
        case 11:
          return TyreCompounds.SUPERSOFT
        case 12:
          return TyreCompounds.SOFT
        case 13:
          return TyreCompounds.MEDIUM
        case 14:
          return TyreCompounds.HARD
        case 15:
          return TyreCompounds.WET
        case 16:
          return TyreCompounds.C5
        case 17:
          return TyreCompounds.C4
        case 18:
          return TyreCompounds.C3
        case 19:
          return TyreCompounds.C2
        case 20:
          return TyreCompounds.C1
        case 21:
          return TyreCompounds.C0
        default:
          return TyreCompounds.UNKNOWN
      }
    },
  },
})
