import type { WsConnectionStatus } from '@/types/store/wsConnectionStatus'
import { defineStore } from 'pinia'

export const useWsStore = defineStore('wsStore', {
  state: () =>
    ({
      wsConnected: false,
    }) as WsConnectionStatus,
  actions: {
    /**
     * Set the WebSocket connection status.
     *
     * @param {boolean} connected - The new WebSocket connection status.
     */
    setWsConnected(connected: boolean) {
      this.wsConnected = connected
    },
  },
})
