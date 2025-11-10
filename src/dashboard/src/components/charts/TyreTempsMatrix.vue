<template>
  <chart-wrapper title="Tyre Temperatures (°C)"
    icon="<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><path d='M12 14v6'></path><path d='M12 2a10 10 0 0 0-7.32 3.25'></path><path d='M12 2a10 10 0 0 1 7.32 3.25'></path><path d='M21 9a10 10 0 0 1-9 13 10 10 0 0 1-9-13'></path><path d='M3 9a10 10 0 0 1 9-7 10 10 0 0 1 9 7'></path></svg>">
    <div class="grid grid-cols-4 grid-rows-4 gap-2 h-full w-full">
      <!-- Front Left Tyre (Left Column, Top Row) -->
      <div class="col-span-2 row-span-2 flex flex-col items-center justify-center">
        <div class="text-sm text-gray-400 mb-2 w-full text-left pl-4">Front Left</div>
        <div class="flex flex-row items-center w-full">
          <div class="flex flex-col items-end mr-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[2], props.data?.compound) }">
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 mr-3 inner-tyre-icon" />
              {{ props.data?.innerTemp?.[2]?.toFixed(1) ?? '0.0' }}°C
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[2], props.data?.compound) }">
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 mr-3" />
              {{ props.data?.outerTemp?.[2]?.toFixed(1) ?? '0.0' }}°C
            </div>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-2"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[2], props.data?.compound)" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[2], props.data?.compound)" />
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
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[3], props.data?.compound)" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[3], props.data?.compound)" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
          <div class="flex flex-col items-start ml-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[3], props.data?.compound) }">
              {{ props.data?.innerTemp?.[3]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 ml-3 inner-tyre-icon" />
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[3], props.data?.compound) }">
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
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[0], props.data?.compound) }">
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 mr-3 inner-tyre-icon" />
              {{ props.data?.innerTemp?.[0]?.toFixed(1) ?? '0.0' }}°C
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[0], props.data?.compound) }">
              <img src="/assets/icons/tyres/OuterTyre_Icon.svg" alt="Outer Tyre Temp Icon" class="w-3 h-3 mr-3" />
              {{ props.data?.outerTemp?.[0]?.toFixed(1) ?? '0.0' }}°C
            </div>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="70" height="100" viewBox="0 0 70 100" id="tyre-0"
            class="w-12 h-12">
            <defs>
              <radialGradient id="outerGradient" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[0], props.data?.compound)" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[0], props.data?.compound)" />
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
                <stop offset="45%" :stop-color="getTempColor(props.data?.innerTemp?.[1], props.data?.compound)" />
                <stop offset="100%" :stop-color="getTempColor(props.data?.outerTemp?.[1], props.data?.compound)" />
              </radialGradient>
            </defs>
            <g id="OuterTyre">
              <rect x="0" y="0" width="70" height="100" rx="15" ry="15" fill="url(#outerGradient)" opacity="1" />
            </g>
          </svg>
          <div class="flex flex-col items-start ml-4">
            <div class="flex items-center text-xl font-bold" :style="{ color: getTempColor(props.data?.innerTemp?.[1], props.data?.compound) }">
              {{ props.data?.innerTemp?.[1]?.toFixed(1) ?? '0.0' }}°C
              <img src="/assets/icons/tyres/InnerTyre_Icon.svg" alt="Inner Tyre Temp Icon"
                class="w-4 h-4 ml-3 inner-tyre-icon" />
            </div>
            <div class="flex items-center text-md font-semibold" :style="{ color: getTempColor(props.data?.outerTemp?.[1], props.data?.compound) }">
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
import { TyreCompounds } from '../../constants.ts'

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
    compound?: number
  }
}>()

function getTempColor(temp: number | undefined, compound: number | undefined): string {
  /**
   * Determines the color code based on tyre temperature.
   *
   * @param temp - The tyre temperature.
   * @param compound - The tyre compound type.
   * @returns The corresponding color code as a string.
   */
  const compoundType = CompoundIdToString(compound ?? 0)

  let coldBound = 80
  let coldWarmBound = 90
  let warmHotBound = 105
  let hotBound = 115

  switch (compoundType) {
    case TyreCompounds.INTERMEDIATE:
      coldBound = 60
      coldWarmBound = 75
      warmHotBound = 95
      hotBound = 100
      break
    case TyreCompounds.WET:
      coldBound = 50
      coldWarmBound = 60
      warmHotBound = 85
      hotBound = 90
      break
    case TyreCompounds.DRY:
      // Currently no specific bounds for this compound, using default values
      break
    case TyreCompounds.SUPERSOFT:
      // Currently no specific bounds for this compound, using default values
      break
    case TyreCompounds.SOFT:
      // Currently no specific bounds for this compound, using default values
      break
    case TyreCompounds.MEDIUM:
      // Currently no specific bounds for this compound, using default values
      break
    case TyreCompounds.HARD:
      // Currently no specific bounds for this compound, using default values
      break
    case TyreCompounds.C5:
      coldBound = 75
      coldWarmBound = 85
      warmHotBound = 105
      hotBound = 110
      break
    case TyreCompounds.C4:
      coldBound = 77.5
      coldWarmBound = 87.5
      warmHotBound = 107.5
      hotBound = 117.5
      break
    case TyreCompounds.C3:
      coldBound = 80
      coldWarmBound = 90
      warmHotBound = 110
      hotBound = 120
      break
    case TyreCompounds.C2:
      coldBound = 82.5
      coldWarmBound = 92.5
      warmHotBound = 112.5
      hotBound = 122.5
      break
    case TyreCompounds.C1:
    case TyreCompounds.C0:
      coldBound = 85
      coldWarmBound = 95
      warmHotBound = 115
      hotBound = 125
      break
  }

  if (!temp || temp < coldBound) return tempColorCodes.cold
  else if (temp >= coldBound && temp < coldWarmBound) return tempColorCodes.coldWarm
  else if (temp >= coldWarmBound && temp < warmHotBound) return tempColorCodes.warm
  else if (temp >= warmHotBound && temp < hotBound) return tempColorCodes.warmHot
  else {
    return tempColorCodes.hot
  }
}

function CompoundIdToString(compoundId: number | undefined): TyreCompounds {
  /**
   * Maps the compound ID from telemetry data to TyreCompounds enum according to F1 2023 telemetry spec.
   *
   * @param compoundId - The compound ID from telemetry data.
   * @returns The corresponding TyreCompounds enum value.
   */
  switch (compoundId) {
    case 7:
      return TyreCompounds.INTERMEDIATE
    case 8:
      return TyreCompounds.WET
    case 9:
      return TyreCompounds.DRY
    case 10:
      return TyreCompounds.WET
    case 11:
      return TyreCompounds.SUPERSOFT
    case 12:
      return TyreCompounds.SOFT
    case 13:
      return TyreCompounds.MEDIUM
    case 14:
      return TyreCompounds.HARD
    case 15:
      return TyreCompounds.WET
    case 16:
      return TyreCompounds.C5
    case 17:
      return TyreCompounds.C4
    case 18:
      return TyreCompounds.C3
    case 19:
      return TyreCompounds.C2
    case 20:
      return TyreCompounds.C1
    case 21:
      return TyreCompounds.C0
    default:
      return TyreCompounds.UNKNOWN
  }
}
</script>
