<template>
  <div id="app" class="p-4">
    <header class="mb-4 grid grid-cols-3 items-center bg-gray-800 p-4 rounded-lg shadow-lg">
      <div class="flex items-center">
        <span
          class="w-3 h-3 rounded-full mr-3 relative group"
          :class="wsStore.wsConnected ? 'bg-green-500' : 'bg-red-500 animate-pulse'"
          style="box-shadow: 0 0 8px 2px currentColor"
          aria-label="Connection status"
          :title="wsStore.wsConnected ? 'Backend Connected' : 'Backend not Connected'"
        >
        </span>
        <img
          v-if="wakeLock.isSupported"
          :src="iconSrc"
          @click="toggleWakeLock"
          alt="Wake-Lock-Icon"
          class="w-6 h-6 ml-2"
          :title="wakeLock.isActive ? 'Disable Wake Lock' : 'Enable Wake Lock'"
        />
      </div>
      <h1 class="text-2xl font-bold text-gray-100 flex items-center mb-2 sm:mb-0 justify-center">
        <img src="/assets/icons/Dashboard_Icon.svg" alt="Dashboard Icon" class="w-7 h-7 mr-3" />
        F1 Telemetry Dashboard
      </h1>

      <!-- Header to toggle widgets -->
      <div class="flex flex-wrap gap-2 justify-end">
        <button
          v-for="widget in widgets"
          :key="widget.id"
          @click="toggleWidget(widget.id)"
          :class="`flex items-center px-4 py-2 rounded-full text-sm font-medium transition-colors duration-200
            ${widget.visible ? 'bg-indigo-600 hover:bg-indigo-700 text-white' : 'bg-gray-700 hover:bg-gray-600 text-gray-300'}`"
        >
          <span class="mr-2 lucide">
            <img v-if="widget.icon" :src="widget.icon" class="w-4 h-4" />
          </span>
          <span>{{ widget.name }}</span>
        </button>
      </div>
    </header>
    <!-- Main Dashboard Grid Layout -->
    <div class="flex flex-wrap gap-6">
      <template v-for="widget in widgets" :key="widget.id">
        <div v-if="widget.visible" class="w-fit h-fit">
          <component :is="widget.component" :data="getWidgetData(widget.id)"></component>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped></style>

<script setup lang="ts">
import { ref, markRaw, computed } from 'vue'
import { useWakeLock } from '@vueuse/core'
import TyreTemps from '@/components/charts/TyreTempsWidget.vue'
import { useTyreStore } from '@/stores/TyreStore'
import { useWsStore } from '@/stores/WsStore'

const tyreStore = useTyreStore()
const wsStore = useWsStore()

const wakeLock = ref(useWakeLock())

const widgets = ref([
  {
    id: 'TyreTemps',
    name: 'Tyre Temps',
    icon: '/assets/icons/Tyre_Icon.svg',
    component: markRaw(TyreTemps),
    visible: true,
  },
  // Add more widgets as needed
])

const iconSrc = computed(() => {
  const iconName = wakeLock.value.isActive ? 'Visible_Icon.svg' : 'Invisible_Icon.svg'
  return new URL(`/assets/icons/${iconName}`, import.meta.url).href
})

/**
 * Fetch data based on the widget ID.
 *
 * @param widgetId - The ID of the widget to fetch data for.
 * @returns The data fetched for the widget.
 */
function getWidgetData(widgetId: string) {
  if (widgetId === 'TyreTemps') {
    return {
      innerTemp: tyreStore.innerTyreTemps,
      outerTemp: tyreStore.outerTyreTemps,
      compound: tyreStore.tyreCompound,
    }
  }

  return {}
}

/**
 * Toggle the visibility of a widget.
 *
 * @param widgetId - The ID of the widget to toggle.
 */
function toggleWidget(widgetId: string) {
  const widget = widgets.value.find((w) => w.id === widgetId)
  if (widget) {
    widget.visible = !widget.visible
  }
}

/**
 * Toggles the wake lock.
 */
function toggleWakeLock() {
  if (wakeLock.value.isActive) {
    wakeLock.value.release()
  } else {
    wakeLock.value.request('screen')
  }
}
</script>
