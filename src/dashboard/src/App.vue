<template>
  <div id="app" class="p-4">
    <header
      class="mb-4 flex flex-col sm:flex-row justify-between items-center bg-gray-800 p-4 rounded-lg shadow-lg"
    >
      <span
        class="w-3 h-3 rounded-full mr-3 relative group"
        :class="wsStore.wsConnected ? 'bg-green-500' : 'bg-red-500 animate-pulse'"
        style="box-shadow: 0 0 8px 2px currentColor;"
        aria-label="Connection status"
      >
        <span
          class="absolute left-1/2 bottom-full mb-2 px-2 py-1 rounded bg-gray-900 text-gray-100 text-xs whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-10"
        >
          {{ wsStore.wsConnected ? 'Backend Connected' : 'Backend not Connected' }}
        </span>
      </span>
      <h1 class="text-2xl font-bold text-gray-100 flex items-center mb-2 sm:mb-0">
        <img src="/assets/icons/Dashboard_Icon.svg" alt="Dashboard Icon" class="w-7 h-7 mr-3" />
        F1 Telemetry Dashboard
      </h1>

      <!-- Header to toggle widgets -->
      <div class="flex flex-wrap gap-2">
        <button
          v-for="widget in widgets"
          :key="widget.id"
          @click="toggleWidget(widget.id)"
          :class="`flex items-center px-4 py-2 rounded-full text-sm font-medium transition-colors duration-200 
            ${widget.visible ? 'bg-indigo-600 hover:bg-indigo-700 text-white' : 'bg-gray-700 hover:bg-gray-600 text-gray-300'}`"
        >
          <span v-html="widget.icon" class="mr-2 lucide"></span>
          <span>{{ widget.name }}</span>
        </button>
      </div>
    </header>
    <!-- Main Dashboard Grid Layout -->
    <div class="flex flex-wrap gap-6">
      <template v-for="widget in widgets" :key="widget.id">
      <div v-if="widget.visible" class="w-fit h-fit">
        <component 
        :is="widget.component" 
        :data="getWidgetData(widget.id)"
        ></component>
      </div>
      </template>
    </div>
    </div>
  </template>

  <style scoped></style>

  <script setup lang="ts">
  import { ref, markRaw } from 'vue'
  import TyreInfo from './components/charts/TyreTempsWidget.vue'
  import { useTyreStore } from './stores/TyreStore'
  import { useWsStore } from './stores/WsStore'

  const tyreStore = useTyreStore()
  const wsStore = useWsStore()

  const widgets = ref([
    {
    id: 'TyreInfo',
    name: 'Tyre Info',
    icon: `<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><path d='M12 14v6'></path><path d='M12 2a10 10 0 0 0-7.32 3.25'></path><path d='M12 2a10 10 0 0 1 7.32 3.25'></path><path d='M21 9a10 10 0 0 1-9 13 10 10 0 0 1-9-13'></path><path d='M3 9a10 10 0 0 1 9-7 10 10 0 0 1 9 7'></path></svg>`,
    component: markRaw(TyreInfo),
    visible: true
    },
    // Add more widgets as needed
  ])

  function toggleWidget(widgetId: string) {
  const widget = widgets.value.find((w) => w.id === widgetId)
  if (widget) {
    widget.visible = !widget.visible
  }
}

function getWidgetData(widgetId: string) {
  if (widgetId === 'TyreInfo') {
    return { 'innerTemp': tyreStore.innerTyreTemps, 'outerTemp': tyreStore.outerTyreTemps, 'compound': tyreStore.actualTyreCompoundId }
  }
  
  return {}
}
</script>
