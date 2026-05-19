import { createApp } from 'vue'
import App from '@/App.vue'
import { createPinia } from 'pinia'
import { useWsStore } from '@/stores/WsStore'
import { useTyreStore } from '@/stores/TyreStore'
import { parserMap } from '@/services/parsers/parserMapper'

createApp(App).use(createPinia()).mount('#app')

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string
    }
  }
}

// Setup stores
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
  // Parse message data as JSON
  let message
  try {
    message = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  // Extract data and header
  const data = message?.['data']
  const header = data?.['M_header']
  if (!data) {
    console.warn('Received WS message with missing data:', message)
    return
  }
  if (!header) {
    console.warn('Received WS message with missing header:', message)
    return
  }

  // Get player car index in array
  const player_id = header['M_playerCarIndex']
  const parser = parserMap[header['M_packetFormat']]

  switch (message.topic) {
    case 'telemetry.car_telemetry': {

      const carData = parser.parseCarTelemetry(data, player_id)

      tyreStore.updateTyreTemps(carData.innerTyreTemps, carData.outerTyreTemps)
      break
    }
    case 'telemetry.car_status': {
      // Handle car status data
      const statusData = parser.parseCarStatusTelemetry(data, player_id)
      tyreStore.updateTyreCompound(statusData.actualTyreCompoundId)
      break
    }
    default:
      // ignore other topics
      break
  }
}
