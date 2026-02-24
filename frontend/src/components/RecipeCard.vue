<script setup lang="ts">
import type { Recipe } from '../types'

defineProps<{
  recipe: Recipe
}>()
</script>

<template>
  <router-link
    :to="`/recipes/${recipe.id}`"
    class="block bg-white rounded-xl border border-stone-200 overflow-hidden hover:border-orange-200 hover:shadow-sm transition-all"
  >
    <img
      v-if="recipe.cover_image"
      :src="recipe.cover_image"
      :alt="recipe.name"
      class="w-full h-36 object-cover"
    />
    <div v-else class="w-full h-36 bg-stone-100 flex items-center justify-center">
      <span class="text-3xl">🍽️</span>
    </div>
    <div class="p-3">
      <h3 class="font-medium text-stone-800 text-sm">{{ recipe.name }}</h3>
      <div class="flex gap-1.5 mt-2 flex-wrap">
        <span v-if="recipe.difficulty" class="px-1.5 py-0.5 bg-stone-100 text-stone-500 text-xs rounded">
          {{ recipe.difficulty }}
        </span>
        <span v-if="recipe.cook_time" class="px-1.5 py-0.5 bg-stone-100 text-stone-500 text-xs rounded">
          {{ recipe.cook_time }}分钟
        </span>
        <span v-for="tag in recipe.tags.slice(0, 3)" :key="tag" class="px-1.5 py-0.5 bg-orange-50 text-orange-600 text-xs rounded">
          {{ tag }}
        </span>
      </div>
      <div v-if="recipe.ingredients.length" class="mt-2 text-xs text-stone-400 truncate">
        {{ recipe.ingredients.map(i => i.name).join('、') }}
      </div>
    </div>
  </router-link>
</template>
