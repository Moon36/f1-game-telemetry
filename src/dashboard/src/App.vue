<template>
  <div id="app" class="p-4">
    <header
      class="mb-4 flex flex-col sm:flex-row justify-between items-center bg-gray-800 p-4 rounded-lg shadow-lg"
    >
      <h1 class="text-2xl font-bold text-gray-100 flex items-center mb-2 sm:mb-0">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="28"
          height="28"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="lucide mr-2"
        >
          <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
          <line x1="3" x2="21" y1="9" y2="9"></line>
          <line x1="3" x2="21" y1="15" y2="15"></line>
          <line x1="9" x2="9" y1="3" y2="21"></line>
          <line x1="15" x2="15" y1="3" y2="21"></line>
        </svg>
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
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 auto-rows-fr">
      <template v-for="widget in widgets" :key="widget.id">
        <component v-if="widget.visible" :is="widget.component" :data="widget.data"></component>
      </template>
    </div>
  </div>
</template>

<style scoped></style>

<script setup lang="ts">
import { ref, markRaw } from 'vue'
import TyreInfo from './components/charts/TyreTempsMatrix.vue'

const widgets = ref([
  {
    id: 'TyreInfo',
    name: 'Tyre Info',
    icon: `<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><path d='M12 14v6'></path><path d='M12 2a10 10 0 0 0-7.32 3.25'></path><path d='M12 2a10 10 0 0 1 7.32 3.25'></path><path d='M21 9a10 10 0 0 1-9 13 10 10 0 0 1-9-13'></path><path d='M3 9a10 10 0 0 1 9-7 10 10 0 0 1 9 7'></path></svg>`,
    component: markRaw(TyreInfo),
    visible: true,
    data: [0, 0, 0, 0],
  },
  // Add more widgets as needed
])

function toggleWidget(widgetId: string) {
  const widget = widgets.value.find((w) => w.id === widgetId)
  if (widget) {
    widget.visible = !widget.visible
  }
}
</script>
