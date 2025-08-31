<template>
  <div
    class="rounded-lg bg-gray-800 border border-gray-700 shadow-md p-4 transition-all duration-300 h-full w-full flex flex-col"
  >
    <div class="flex items-center text-gray-300 mb-2">
      <span class="mr-2 lucide" v-html="icon"></span>
      <h3 class="text-sm font-semibold">{{ title }}</h3>
    </div>
    <div class="flex-grow">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

defineProps(['title', 'icon'])
</script>
