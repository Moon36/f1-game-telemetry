import { reactive } from 'vue'

export const store = reactive({
  wsConnected: false,
  tyreTemps: [0, 0, 0, 0] as number[],
})
