<template>
  <chart-wrapper title="Tyre Temperatures (°C)"
    icon="<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><path d='M12 14v6'></path><path d='M12 2a10 10 0 0 0-7.32 3.25'></path><path d='M12 2a10 10 0 0 1 7.32 3.25'></path><path d='M21 9a10 10 0 0 1-9 13 10 10 0 0 1-9-13'></path><path d='M3 9a10 10 0 0 1 9-7 10 10 0 0 1 9 7'></path></svg>">
    <div class="grid grid-cols-4 grid-rows-4 gap-2 h-full w-full">
      <!-- Front Left Tyre (Left Column, Top Row) -->
      <div class="col-span-2 row-span-2 flex flex-col items-center justify-center">
        <div class="text-sm text-gray-400 mb-2 w-full text-left pl-4">Front Left</div>
        <div class="flex flex-row items-center w-full">
          <div class="flex flex-col items-end mr-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[2]) }">
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 mr-3 inner-tyre-icon" />
              {{ props.data?.innerTemp?.[2]?.toFixed(1) ?? '0.0' }}°C
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[2]) }">
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 mr-3" />
              {{ props.data?.outerTemp?.[2]?.toFixed(1) ?? '0.0' }}°C
            </div>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-2"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[2])" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[2])" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
        </div>
      </div>

      <!-- Front Right Tyre (Right Column, Top Row) -->
      <div class="col-span-2 row-span-2 flex flex-col items-center justify-center">
        <div class="text-sm text-gray-400 mb-2 w-full text-right pr-4">Front Right</div>
        <div class="flex flex-row items-center w-full justify-end">
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-3"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[3])" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[3])" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
          <div class="flex flex-col items-start ml-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[3]) }">
              {{ props.data?.innerTemp?.[3]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 ml-3 inner-tyre-icon" />
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[3]) }">
              {{ props.data?.outerTemp?.[3]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 ml-3" />
            </div>
          </div>
        </div>
      </div>

      <!-- Rear Left Tyre (Left Column, Bottom Row) -->
      <div class="col-span-2 row-span-2 flex flex-col items-center justify-center">
        <div class="flex flex-row items-center w-full">
          <div class="flex flex-col items-end mr-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[0]) }">
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 mr-3 inner-tyre-icon" />
              {{ props.data?.innerTemp?.[0]?.toFixed(1) ?? '0.0' }}°C
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[0]) }">
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 mr-3" />
              {{ props.data?.outerTemp?.[0]?.toFixed(1) ?? '0.0' }}°C
            </div>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-0"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[0])" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[0])" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
        </div>
        <div class="text-sm text-gray-400 mt-2 w-full text-left pl-4">Rear Left</div>
      </div>

      <!-- Rear Right Tyre (Right Column, Bottom Row) -->
      <div class="col-span-2 row-span-2 flex flex-col items-center justify-center">
        <div class="flex flex-row items-center w-full justify-end">
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-1"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[1])" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[1])" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
          <div class="flex flex-col items-start ml-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[1]) }">
              {{ props.data?.innerTemp?.[1]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 ml-3 inner-tyre-icon" />
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[1]) }">
              {{ props.data?.outerTemp?.[1]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 ml-3" />
            </div>
          </div>
        </div>
        <div class="text-sm text-gray-400 mt-2 w-full text-right pr-4">Rear Right</div>
      </div>
    </div>
  </chart-wrapper>
</template>

<script setup lang="ts">
import ChartWrapper from '../ChartWrapper.vue'

const tempColorCodes = {
  cold: '#3b82f6',
  coldWarm: '#86efac',
  warm: '#22c55e',
  warmHot: '#f97316',
  hot: '#ef4444',
}

const props = defineProps<{
  data?: {
    innerTemp?: number[]
    outerTemp?: number[]
  }
}>()

// TODO: Make this dependent on tyre compound and add intermediate colors
function getTempColor(temp: number | undefined): string {
  if (!temp || temp < 90) return tempColorCodes.cold
  else if (temp >= 90 && temp < 100) return tempColorCodes.coldWarm
  else if (temp >= 100 && temp < 110) return tempColorCodes.warm
  else if (temp >= 110 && temp < 120) return tempColorCodes.warmHot
  else {
    return tempColorCodes.hot
  }
}
</script>
