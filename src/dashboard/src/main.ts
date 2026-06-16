import { createApp } from 'vue'
import App from '@/App.vue'
import { createPinia } from 'pinia'
import { useWsStore } from '@/stores/WsStore'
import { useTyreStore } from '@/stores/TyreStore'
import { parserMap } from '@/services/parsers/parserMapper'
import * as constants from '@/utils/constants.ts'
import type { GenericRawTelemetry, Header } from './types/packets/packetDefinitions'

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

// Create WebSocket connection
connectWebSocket(wsPort, handleWSMessage)

/**
 * Creates a WebSocket connection to a server on localhost and the specified port.
 * The provided messageHandlerCB callback function is called whenever a new message is received.
 * This method reconnects automatically in the given time interval, if connection is lost (default is 1000ms/1s).
 *
 * @param port - The port number to connect to.
 * @param messageHandlerCB - The callback function to handle incoming messages.
 * @param reconnectInterval - The time interval in milliseconds to wait before attempting a reconnection.
 */
function connectWebSocket(
  port: string,
  messageHandlerCB: (event: MessageEvent) => void,
  reconnectInterval: number = 1000,
) {
  console.log(`Trying to connect to WebSocket on port ${port}...`)
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    console.log('WebSocket connected')
    wsStore.setWsConnected(true)
  }

  ws.onmessage = (event) => {
    messageHandlerCB(event)
  }

  ws.onclose = function (event) {
    console.log('Socket was closed.', event.reason)
    wsStore.setWsConnected(false)
    setTimeout(function () {
      connectWebSocket(port, messageHandlerCB, reconnectInterval)
    }, reconnectInterval)
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    ws.close()
    wsStore.setWsConnected(false)
  }
}

/**
 * Handles messages from the WebSocket.
 *
 * @param event - The WebSocket message event.
 */
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
  const packet: GenericRawTelemetry = message?.data
  const header: Header = packet?.M_header
  if (!packet) {
    console.warn('Received WS message with missing data:', message)
    return
  }
  if (!header) {
    console.warn('Received WS message with missing header:', message)
    return
  }

  // Get player car index in array
  const player_id = header.M_playerCarIndex
  const parser = parserMap[header.M_packetFormat]

  switch (message.topic) {
    case constants.TOPIC_CAR_TELEMETRY_DATA: {
      const carData = parser.parseCarTelemetry(packet, player_id)

      tyreStore.updateTyreTemps(carData.innerTyreTemps, carData.outerTyreTemps)
      break
    }
    case constants.TOPIC_CAR_STATUS_DATA: {
      // Handle car status data
      const statusData = parser.parseCarStatusTelemetry(packet, player_id)
      tyreStore.updateTyreCompound(statusData.actualTyreCompoundId)
      break
    }
    default:
      // ignore other topics
      break
  }
}
