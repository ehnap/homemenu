<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getSharedMealPlan } from '../api/mealPlan'
import type { MealPlan, ShoppingList, MealPlanItem, Recipe } from '../types'
import { MEAL_TYPE_LABELS, WEEKDAY_LABELS } from '../types'

const route = useRoute()
const token = route.params.token as string

const plan = ref<MealPlan | null>(null)
const shoppingList = ref<ShoppingList | null>(null)
const loading = ref(true)
const error = ref('')
const selectedRecipe = ref<Recipe | null>(null)

onMounted(async () => {
  try {
    const data = await getSharedMealPlan(token)
    plan.value = data.meal_plan
    shoppingList.value = data.shopping_list
  } catch (e: any) {
    error.value = '分享链接无效或已过期'
  } finally {
    loading.value = false
  }
})

const dates = computed(() => {
  if (!plan.value) return []
  const result: string[] = []
  const start = new Date(plan.value.start_date)
  const end = new Date(plan.value.end_date)
  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    result.push(d.toISOString().split('T')[0]!)
  }
  return result
})

const mealTypes = computed(() => {
  if (!plan.value?.config?.meal_types?.length) return ['lunch', 'dinner']
  return plan.value.config.meal_types
})

function getItems(date: string, mealType: string): MealPlanItem[] {
  return (plan.value?.items || [])
    .filter(i => i.date === date && i.meal_type === mealType)
    .sort((a, b) => a.sort_order - b.sort_order)
}

function getWeekday(date: string) {
  const d = new Date(date)
  const day = d.getDay()
  return WEEKDAY_LABELS[day === 0 ? 6 : day - 1]
}

function formatDate(date: string) {
  return date.substring(5).replace('-', '/')
}

function showRecipe(recipe?: Recipe) {
  selectedRecipe.value = recipe || null
}
</script>

<template>
  <div class="min-h-screen bg-stone-50">
    <div class="max-w-7xl mx-auto px-4 py-6">
      <div v-if="loading" class="text-center py-12 text-stone-500">加载中...</div>

      <div v-else-if="error" class="text-center py-12">
        <p class="text-stone-500">{{ error }}</p>
      </div>

      <template v-else-if="plan">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-orange-600">家味</h1>
          <h2 class="text-lg text-stone-700 mt-1">{{ plan.name }}</h2>
          <p class="text-sm text-stone-500">{{ formatDate(plan.start_date) }} - {{ formatDate(plan.end_date) }}</p>
        </div>

        <!-- Meal plan grid -->
        <div class="overflow-x-auto mb-8">
          <table class="w-full bg-white rounded-xl border border-stone-200">
            <thead>
              <tr>
                <th class="px-3 py-2 text-sm text-stone-500 font-medium"></th>
                <th v-for="date in dates" :key="date" class="px-3 py-2 text-center">
                  <div class="text-sm font-medium text-stone-800">{{ getWeekday(date) }}</div>
                  <div class="text-xs text-stone-500">{{ formatDate(date) }}</div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mealType in mealTypes" :key="mealType" class="border-t border-stone-100">
                <td class="px-3 py-3 text-sm font-medium text-stone-600">{{ MEAL_TYPE_LABELS[mealType] }}</td>
                <td v-for="date in dates" :key="date" class="px-2 py-2">
                  <div class="space-y-1">
                    <button
                      v-for="item in getItems(date, mealType)"
                      :key="item.id"
                      @click="showRecipe(item.recipe)"
                      class="block w-full text-left px-2 py-1 bg-orange-50 text-orange-700 text-xs rounded hover:bg-orange-100 transition-colors"
                    >
                      {{ item.recipe?.name || '未知菜品' }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Shopping list -->
        <div v-if="shoppingList?.daily" class="max-w-2xl mx-auto">
          <h2 class="text-lg font-semibold text-stone-800 mb-4">购物清单</h2>
          <div class="space-y-3">
            <div v-for="day in shoppingList.daily" :key="day.date" class="bg-white rounded-xl border border-stone-200">
              <div class="px-4 py-2 border-b border-stone-100">
                <span class="text-sm font-medium text-stone-700">{{ formatDate(day.date) }}</span>
              </div>
              <div class="px-4 py-2">
                <div v-for="item in day.items" :key="item.name" class="flex justify-between py-1 text-sm">
                  <span class="text-stone-700">{{ item.name }}</span>
                  <span class="text-stone-500">{{ item.amount }}{{ item.unit }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- Recipe detail modal -->
      <div v-if="selectedRecipe" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="selectedRecipe = null">
        <div class="bg-white rounded-xl p-6 max-w-lg w-full mx-4 max-h-[80vh] overflow-y-auto">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-stone-800">{{ selectedRecipe.name }}</h3>
            <button @click="selectedRecipe = null" class="text-stone-400 hover:text-stone-600 text-xl">&times;</button>
          </div>

          <div v-if="selectedRecipe.ingredients?.length" class="mb-4">
            <h4 class="text-sm font-medium text-stone-600 mb-2">食材</h4>
            <div v-for="ing in selectedRecipe.ingredients" :key="ing.name" class="flex justify-between text-sm py-1">
              <span>{{ ing.name }}</span>
              <span class="text-stone-500">{{ ing.amount }}{{ ing.unit }}</span>
            </div>
          </div>

          <div v-if="selectedRecipe.steps?.length">
            <h4 class="text-sm font-medium text-stone-600 mb-2">做法</h4>
            <ol class="space-y-2">
              <li v-for="step in selectedRecipe.steps" :key="step.order" class="flex gap-2 text-sm">
                <span class="text-orange-600 font-medium">{{ step.order }}.</span>
                <span class="text-stone-700">{{ step.description }}</span>
              </li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
