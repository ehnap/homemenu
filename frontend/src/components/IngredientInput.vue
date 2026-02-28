<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Ingredient } from '../types'
import { suggestIngredients } from '../api/recipe'

const model = defineModel<Ingredient>({ required: true })

defineProps<{
  canRemove: boolean
}>()

const emit = defineEmits<{
  remove: []
}>()

const suggestions = ref<string[]>([])
const showSuggestions = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(() => model.value.name, (val) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (!val || val.trim().length === 0) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }
  debounceTimer = setTimeout(async () => {
    try {
      const result = await suggestIngredients(val.trim())
      suggestions.value = result
      showSuggestions.value = result.length > 0
    } catch {
      suggestions.value = []
      showSuggestions.value = false
    }
  }, 300)
})

function selectSuggestion(name: string) {
  model.value.name = name
  suggestions.value = []
  showSuggestions.value = false
}

function onBlur() {
  // Delay to allow click event on suggestion item
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}

function onFocus() {
  if (suggestions.value.length > 0) {
    showSuggestions.value = true
  }
}
</script>

<template>
  <div class="flex flex-wrap gap-2 items-center">
    <div class="relative min-w-0 flex-1 basis-full sm:basis-0">
      <input
        v-model="model.name"
        type="text"
        placeholder="食材名"
        @blur="onBlur"
        @focus="onFocus"
        class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
      />
      <ul
        v-if="showSuggestions"
        class="absolute z-10 left-0 right-0 mt-1 bg-white border border-stone-200 rounded-lg shadow-lg max-h-40 overflow-y-auto"
      >
        <li
          v-for="name in suggestions"
          :key="name"
          @mousedown.prevent="selectSuggestion(name)"
          class="px-3 py-2 text-sm text-stone-700 hover:bg-orange-50 hover:text-orange-700 cursor-pointer"
        >{{ name }}</li>
      </ul>
    </div>
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
