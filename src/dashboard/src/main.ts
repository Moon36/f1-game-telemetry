import { createApp } from 'vue'
import App from './App.vue'
import { store } from './store'

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string
    }
  }
}

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? '8282'
const reconnectInterval = 1000


// Create WebSocket connection
connectWebSocket(wsPort)

function connectWebSocket(port: string) {
  console.log(`Trying to connect to WebSocket on port ${port}...`)
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    console.log('WebSocket connected')
    store.wsConnected = true
  }

  ws.onmessage = (event) => {
    handleWSMessage(event)
  };

  ws.onclose = function (event) {
    console.log('Socket was closed.', event.reason);
    store.wsConnected = false
    setTimeout(function () {
      connectWebSocket(port);
    }, reconnectInterval);
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    ws.close()
    store.wsConnected = false
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


createApp(App).mount('#app')
