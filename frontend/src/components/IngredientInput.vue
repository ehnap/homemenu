<script setup lang="ts">
import type { Ingredient } from '../types'

const model = defineModel<Ingredient>({ required: true })

defineProps<{
  canRemove: boolean
}>()

const emit = defineEmits<{
  remove: []
}>()
</script>

<template>
  <div class="flex flex-wrap gap-2 items-center">
    <input
      v-model="model.name"
      type="text"
      placeholder="食材名"
      class="min-w-0 flex-1 basis-full sm:basis-0 px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
    />
    <div class="flex gap-2 items-center flex-1 sm:flex-none">
      <input
        v-model="model.amount"
        type="text"
        placeholder="用量"
        class="min-w-0 flex-1 sm:w-20 sm:flex-none px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
      />
      <input
        v-model="model.unit"
        type="text"
        placeholder="单位"
        class="min-w-0 w-16 shrink-0 px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
      />
      <button
        v-if="canRemove"
        type="button"
        @click="emit('remove')"
        class="text-stone-400 hover:text-red-500 text-lg shrink-0"
      >&times;</button>
      <div v-else class="w-4 shrink-0"></div>
    </div>
  </div>
</template>
