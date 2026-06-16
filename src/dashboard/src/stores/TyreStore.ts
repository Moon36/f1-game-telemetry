import type { TyreTelemetry } from '@/types/store/tyreTelemetry'
import { TYRE_COMPOUNDS } from '@/utils/constants'
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
    tyreCompound(): TYRE_COMPOUNDS {
      switch (this.actualTyreCompoundId) {
        case 7:
          return TYRE_COMPOUNDS.INTERMEDIATE
        case 8:
          return TYRE_COMPOUNDS.WET
        case 9:
          return TYRE_COMPOUNDS.DRY
        case 10:
          return TYRE_COMPOUNDS.WET
        case 11:
          return TYRE_COMPOUNDS.SUPERSOFT
        case 12:
          return TYRE_COMPOUNDS.SOFT
        case 13:
          return TYRE_COMPOUNDS.MEDIUM
        case 14:
          return TYRE_COMPOUNDS.HARD
        case 15:
          return TYRE_COMPOUNDS.WET
        case 16:
          return TYRE_COMPOUNDS.C5
        case 17:
          return TYRE_COMPOUNDS.C4
        case 18:
          return TYRE_COMPOUNDS.C3
        case 19:
          return TYRE_COMPOUNDS.C2
        case 20:
          return TYRE_COMPOUNDS.C1
        case 21:
          return TYRE_COMPOUNDS.C0
        default:
          return TYRE_COMPOUNDS.UNKNOWN
      }
    },
  },
})
