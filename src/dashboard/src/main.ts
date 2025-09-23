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

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? console.error('WebSocket port is not defined')


// Create WebSocket connection
const ws = new WebSocket(`ws://localhost:${wsPort}`)

ws.onopen = () => {
  console.log('WebSocket connected')
}

ws.onmessage = (event) => {
  let data: any
  try {
    data = JSON.parse(event.data)
  } catch (err) {
    console.error('Failed to parse WS message as JSON', err)
    return
  }

  switch (data?.topic) {
    case 'telemetry.car_telemetry': {
      // Handle car telemetry data
      const newTemps = data['data']['M_carTelemetry'][0]['M_tyresInnerTemperature']

      // Mutate array in place to keep reactivity
      if (Array.isArray(newTemps) && newTemps.length === 4) {
        store.tyreTemps.splice(0, 4, ...newTemps)
        console.log('Updated tyre temperatures:', store.tyreTemps)
      } else {
        console.warn('Invalid tyre temperature data:', newTemps)
      }
      break
    }
    default:
      // ignore other topics
      break
  }
}

ws.onerror = (error) => {
  console.error('WebSocket error:', error)
}


createApp(App).mount('#app')
