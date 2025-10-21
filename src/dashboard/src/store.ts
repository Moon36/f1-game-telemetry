import { reactive } from 'vue'

export const store = reactive({
  wsConnected: false,
  innerTyreTemps: [0, 0, 0, 0] as number[],
  outerTyreTemps: [0, 0, 0, 0] as number[],
})
