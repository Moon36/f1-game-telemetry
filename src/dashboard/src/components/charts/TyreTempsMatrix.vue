<template>
  <chart-wrapper
    title="Tyre Temperatures (°C)"
    icon="<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><path d='M12 14v6'></path><path d='M12 2a10 10 0 0 0-7.32 3.25'></path><path d='M12 2a10 10 0 0 1 7.32 3.25'></path><path d='M21 9a10 10 0 0 1-9 13 10 10 0 0 1-9-13'></path><path d='M3 9a10 10 0 0 1 9-7 10 10 0 0 1 9 7'></path></svg>"
  >
    <div class="flex justify-center items-center h-full">
      <div class="grid grid-cols-2 gap-4">
        <div class="flex flex-col items-center">
          <div class="text-2xl font-bold" :style="{ color: getTempColor(data[0]) }">
            {{ data[0].toFixed(1) }}°C
          </div>
          <div class="text-sm text-gray-400">Rear Left</div>
        </div>
        <div class="flex flex-col items-center">
          <div class="text-2xl font-bold" :style="{ color: getTempColor(data[1]) }">
            {{ data[1].toFixed(1) }}°C
          </div>
          <div class="text-sm text-gray-400">Rear Right</div>
        </div>
        <div class="flex flex-col items-center">
          <div class="text-2xl font-bold" :style="{ color: getTempColor(data[2]) }">
            {{ data[2].toFixed(1) }}°C
          </div>
          <div class="text-sm text-gray-400">Front Left</div>
        </div>
        <div class="flex flex-col items-center">
          <div class="text-2xl font-bold" :style="{ color: getTempColor(data[3]) }">
            {{ data[3].toFixed(1) }}°C
          </div>
          <div class="text-sm text-gray-400">Front Right</div>
        </div>
      </div>
    </div>
  </chart-wrapper>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
import ChartWrapper from '../ChartWrapper.vue'

const tempColorCodes = {
  cold: '#3b82f6',
  warm: '#22c55e',
  hot: '#ef4444',
}

const props = defineProps({
  data: {
    type: Array,
    required: true,
    default: () => [0, 0, 0, 0],
  },
})

// TODO: Make this dependent on tyre compound and add intermediate colors
function getTempColor(temp: number) {
  if (temp < 95) return tempColorCodes.cold
  if (temp < 115) return tempColorCodes.warm
  return tempColorCodes.hot
}
</script>
