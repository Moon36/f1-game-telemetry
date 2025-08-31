import { createApp } from 'vue'
import App from './App.vue'

// Access Vite environment variables
const wsPort = import.meta.env.VITE_DASHBOARD_PORT || '8080'

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
