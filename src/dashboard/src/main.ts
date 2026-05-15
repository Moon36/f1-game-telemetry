import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import { useWsStore } from './stores/WsStore'
import { useTyreStore } from './stores/TyreStore'

createApp(App).use(createPinia()).mount('#app')

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string
    }
  }
}

// Access stores
const wsStore = useWsStore()
const tyreStore = useTyreStore()

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? '8282'
const reconnectInterval = 1000

// Create WebSocket connection
connectWebSocket(wsPort)

function connectWebSocket(port: string) {
  console.log(`Trying to connect to WebSocket on port ${port}...`)
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    console.log('WebSocket connected')
    wsStore.setWsConnected(true)
  }

  ws.onmessage = (event) => {
    handleWSMessage(event)
  };

  ws.onclose = function (event) {
    console.log('Socket was closed.', event.reason);
    wsStore.setWsConnected(false)
    setTimeout(function () {
      connectWebSocket(port);
    }, reconnectInterval);
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    ws.close()
    wsStore.setWsConnected(false)
  };
}

function handleWSMessage(event: MessageEvent) {
  let data
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  // Get player car index in array
  const player_id = data['data']['M_header']['M_playerCarIndex']

  switch (data?.topic) {
    case 'telemetry.car_telemetry': {

      // Handle car telemetry data
      const newInnerTemps = data['data']['M_carTelemetry'][player_id]['M_tyresInnerTemperature']
      const newOuterTemps = data['data']['M_carTelemetry'][player_id]['M_tyresSurfaceTemperature']

      // Mutate array in place to keep reactivity
      tyreStore.updateTyreTemps(newInnerTemps, newOuterTemps)
      break
    }
    case 'telemetry.car_status': {
      // Handle car status data
      const compoundId = data['data']['M_carStatusData'][player_id]['M_actualTyreCompound']
      tyreStore.updateTyreCompound(compoundId)
      break
    }
    default:
      // ignore other topics
      break
  }
}
