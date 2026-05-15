import { defineStore } from "pinia"

export const useTyreStore = defineStore("tyreStore", {
  state: () => ({
    innerTyreTemps: [0, 0, 0, 0] as number[],
    outerTyreTemps: [0, 0, 0, 0] as number[],
    actualTyreCompoundId: 0 as number,
  }),
  actions: {
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
    updateTyreCompound(compoundId: number) {
      if (typeof compoundId === 'number') {
        this.actualTyreCompoundId = compoundId
      } else {
        console.warn('Invalid tyre compound ID data:', compoundId)
      }
    }
  }
})