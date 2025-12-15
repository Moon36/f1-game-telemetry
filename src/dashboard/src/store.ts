import { reactive } from 'vue'

export const store = reactive({
  wsConnected: false,                         // Is the WebSocket connection to the backend active
  reWsConnected: false,                       // Is the WebSocket connection to the race engineer active
  
  innerTyreTemps: [0, 0, 0, 0] as number[],   // Inner temperatures of the tyres
  outerTyreTemps: [0, 0, 0, 0] as number[],   // Outer temperatures of the tyres
  actualTyreCompoundId: 0 as number,          // Current tyre compound identifier
})
