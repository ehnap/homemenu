<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getSharedRecipe } from '../api/recipe'
import type { Recipe } from '../types'

const route = useRoute()
const token = route.params.token as string

const recipe = ref<Recipe | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    recipe.value = await getSharedRecipe(token)
  } catch (e: any) {
    error.value = '分享链接无效或已过期'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-stone-50">
    <div class="max-w-3xl mx-auto px-4 py-6">
      <div v-if="loading" class="text-center py-12 text-stone-500">加载中...</div>

      <div v-else-if="error" class="text-center py-12">
        <p class="text-stone-500">{{ error }}</p>
      </div>

      <template v-else-if="recipe">
        <div class="text-center mb-6">
          <h1 class="text-2xl font-bold text-orange-600">家味</h1>
        </div>

        <div class="mb-6">
          <h2 class="text-xl font-bold text-stone-800">{{ recipe.name }}</h2>
          <div class="flex gap-2 mt-2 flex-wrap">
            <span v-if="recipe.difficulty" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
              {{ recipe.difficulty }}
            </span>
            <span v-if="recipe.cook_time" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
              {{ recipe.cook_time }}分钟
            </span>
            <span v-if="recipe.calories" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
              {{ recipe.calories }}卡
            </span>
            <span v-for="tag in recipe.tags" :key="tag" class="px-2 py-0.5 bg-orange-50 text-orange-600 text-xs rounded-full">
              {{ tag }}
            </span>
          </div>
        </div>

        <img
          v-if="recipe.cover_image"
          :src="recipe.cover_image"
          :alt="recipe.name"
          class="w-full h-64 object-cover rounded-xl mb-6"
        />

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="md:col-span-1">
            <div class="bg-white rounded-xl p-4 border border-stone-200">
              <h3 class="font-semibold text-stone-800 mb-3">食材</h3>
              <ul v-if="recipe.ingredients?.length" class="space-y-2">
                <li v-for="ing in recipe.ingredients" :key="ing.name" class="flex justify-between text-sm">
                  <span class="text-stone-700">{{ ing.name }}</span>
                  <span class="text-stone-500">{{ ing.amount }}{{ ing.unit }}</span>
                </li>
              </ul>
              <p v-else class="text-stone-500 text-sm">暂无食材信息</p>
            </div>
          </div>

          <div class="md:col-span-2">
            <div class="bg-white rounded-xl p-4 border border-stone-200">
              <h3 class="font-semibold text-stone-800 mb-3">做法步骤</h3>
              <div v-if="!recipe.steps?.length" class="text-stone-500 text-sm">暂无步骤</div>
              <ol class="space-y-4">
                <li v-for="step in recipe.steps" :key="step.order" class="flex gap-3">
                  <span class="w-6 h-6 bg-orange-100 text-orange-600 rounded-full flex items-center justify-center text-sm font-medium shrink-0">
                    {{ step.order }}
                  </span>
                  <div>
                    <p class="text-sm text-stone-700">{{ step.description }}</p>
                    <img v-if="step.image_url" :src="step.image_url" class="mt-2 rounded-lg max-h-40" />
                  </div>
                </li>
              </ol>
            </div>

            <div v-if="recipe.notes" class="bg-white rounded-xl p-4 border border-stone-200 mt-4">
              <h3 class="font-semibold text-stone-800 mb-2">备注</h3>
              <p class="text-sm text-stone-600">{{ recipe.notes }}</p>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
