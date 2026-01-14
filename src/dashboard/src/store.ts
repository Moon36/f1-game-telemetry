import { reactive } from 'vue'

export const store = reactive({
  isWsConnected: false,                       // Is the WebSocket connection to the backend active
  isReWsConnected: false,                     // Is the WebSocket connection to the race engineer active
  isAudioEnabled: false,                      // Is audio playback enabled (muted/unmuted)
  
  innerTyreTemps: [0, 0, 0, 0] as number[],   // Inner temperatures of the tyres
  outerTyreTemps: [0, 0, 0, 0] as number[],   // Outer temperatures of the tyres
  actualTyreCompoundId: 0 as number,          // Current tyre compound identifier
})
