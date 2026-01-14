import { createApp } from 'vue'
import App from '@/App.vue'
import { store } from '@/store'
import type { AudioChunkMessage } from '@/messages/AudioChunkMessage';
import { audioPlayer } from '@/services/AudioPlayer';

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string,
      RE_BACKEND_PORT: string,
    }
  }
}

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? '8282'
const reWsPort = window.APP_CONFIG?.RE_BACKEND_PORT ?? '8283'
const reconnectInterval = 1000

// Create WebSocket connection
connectWebSocket(wsPort, handleWSMessage, (v) => store.isWsConnected = v)
connectWebSocket(reWsPort, handleREWSMessage, (v) => store.isReWsConnected = v)

function connectWebSocket(port: string, messageHandler: (event: MessageEvent) => void, setConnected?: (connected: boolean) => void) {
  /**
   * Connect to WebSocket server and set up handlers for messages, close, and errors.
   *
   * @param port The port to connect to.
   * @param messageHandler The function to handle incoming messages.
   * @param setConnected Optional function to update connection status.
   */
  console.log(`Trying to connect to WebSocket on port ${port}...`)
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    console.log('WebSocket connected')
    setConnected?.(true)
  }

  ws.onmessage = (event) => {
    messageHandler(event)
  };

  ws.onclose = function (event) {
    console.log('Socket was closed.', event.reason);
    setConnected?.(false)
    setTimeout(function () {
      connectWebSocket(port, messageHandler, setConnected);
    }, reconnectInterval);
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    ws.close()
    setConnected?.(false)
  };
}

function handleWSMessage(event: MessageEvent) {
  /**
   * Handle incoming WebSocket telemetry messages.
   *
   * @param event The WebSocket message event.
   */
  let data
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  switch (data?.topic) {
    case 'telemetry.car_telemetry': {
      // Handle car telemetry data
      const newInnerTemps = data['data']['M_carTelemetry'][0]['M_tyresInnerTemperature']
      const newOuterTemps = data['data']['M_carTelemetry'][0]['M_tyresSurfaceTemperature']

      // Mutate array in place to keep reactivity
      if (Array.isArray(newInnerTemps) && newInnerTemps.length === 4) {
        store.innerTyreTemps.splice(0, 4, ...newInnerTemps)
        console.log('Updated inner tyre temperatures:', store.innerTyreTemps)
      } else {
        console.warn('Invalid inner tyre temperature data:', newInnerTemps)
      }

      if (Array.isArray(newOuterTemps) && newOuterTemps.length === 4) {
        store.outerTyreTemps.splice(0, 4, ...newOuterTemps)
        console.log('Updated outer tyre temperatures:', store.outerTyreTemps)
      } else {
        console.warn('Invalid outer tyre temperature data:', newOuterTemps)
      }
      break
    }
    case 'telemetry.car_status': {
      // Handle car status data
      const compoundId = data['data']['M_carStatusData'][0]['M_actualTyreCompound']
      if (typeof compoundId === 'number') {
        store.actualTyreCompoundId = compoundId
      } else {
        console.warn('Invalid tyre compound ID data:', compoundId)
      }
      break
    }
    default:
      // ignore other topics
      break
  }
}

function handleREWSMessage(event: MessageEvent) {
  /**
   * Handle incoming WebSocket audio chunk messages.
   *
   * @param event The WebSocket message event.
   */
  if (!store.isAudioEnabled) return;

  let data: AudioChunkMessage
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  if (!data.audio_data || data.audio_data.length < 1) return;

  // Handle messages specific to the Race Engineer WebSocket
  console.log('Race Engineer WS message:', data);

  audioPlayer.playAudio(data, true);
}

createApp(App).mount('#app')
