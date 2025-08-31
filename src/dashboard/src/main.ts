import { createApp } from 'vue'
import App from './App.vue'

// Create a runtime config
declare global {
  interface Window {
    APP_CONFIG?: {
      BACKEND_PORT: string
    }
  }
}

const wsPort = window.APP_CONFIG?.BACKEND_PORT ?? console.error('WebSocket port is not defined')
// ToDo: Handle undefined websocket port

console.log(`WebSocket port is set to: ${wsPort}`)

/*
// Create WebSocket connection
const ws = new WebSocket(`ws://localhost:${wsPort}`)

ws.onopen = () => {
  console.log('WebSocket connected')
}

ws.onmessage = (event) => {
  console.log('Received:', event.data)
}

ws.onerror = (error) => {
  console.error('WebSocket error:', error)
}
*/

createApp(App).mount('#app')
