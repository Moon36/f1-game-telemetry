import { defineStore } from "pinia"

export const useWsStore = defineStore("wsStore", {
  state: () => ({
    wsConnected: false,
  }),
  actions: {
    setWsConnected(connected: boolean) {
      this.wsConnected = connected
    }
  }
})