<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMealPlanStore } from '../stores/mealPlan'
import type { PlanConfig } from '../types'

const store = useMealPlanStore()
const router = useRouter()

const today = new Date()
const dayOfWeek = today.getDay()
const mondayOffset = dayOfWeek === 0 ? -6 : 1 - dayOfWeek
const monday = new Date(today)
monday.setDate(today.getDate() + mondayOffset)
const sunday = new Date(monday)
sunday.setDate(monday.getDate() + 6)

function formatDate(d: Date): string {
  return d.toISOString().split('T')[0]!
}

const startDate = ref(formatDate(monday))
const endDate = ref(formatDate(sunday))
const mealTypes = ref(['lunch', 'dinner'])
const tastePreference = ref('')
const preferIngredients = ref('')
const excludeIngredients = ref('')
const loading = ref(false)
const error = ref('')

// Dish counts
const lunchMeat = ref(2)
const lunchVeg = ref(1)
const lunchSoup = ref(1)
const dinnerMeat = ref(1)
const dinnerVeg = ref(1)
const dinnerSoup = ref(0)

function toggleMealType(type: string) {
  const idx = mealTypes.value.indexOf(type)
  if (idx >= 0) {
    mealTypes.value.splice(idx, 1)
  } else {
    mealTypes.value.push(type)
  }
}

async function generate() {
  if (mealTypes.value.length === 0) {
    error.value = '请至少选择一个餐次'
    return
  }

  error.value = ''
  loading.value = true

  const config: PlanConfig = {
    meal_types: mealTypes.value,
    dishes_per_meal: {
      lunch: { meat: lunchMeat.value, vegetable: lunchVeg.value, soup: lunchSoup.value },
      dinner: { meat: dinnerMeat.value, vegetable: dinnerVeg.value, soup: dinnerSoup.value },
      breakfast: { meat: 1, vegetable: 0, soup: 0 },
    },
    taste_preference: tastePreference.value,
    prefer_ingredients: preferIngredients.value ? preferIngredients.value.split(',').map(s => s.trim()) : [],
    exclude_ingredients: excludeIngredients.value ? excludeIngredients.value.split(',').map(s => s.trim()) : [],
  }

  try {
    const plan = await store.generate(startDate.value, endDate.value, config)
    router.push(`/meal-plans/${plan.id}`)
  } catch (e: any) {
    error.value = e.message || '生成失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-stone-800">生成菜单</h1>
      <button @click="router.back()" class="text-sm text-stone-500 hover:text-stone-700">返回</button>
    </div>

    <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg mb-4">{{ error }}</div>

    <form @submit.prevent="generate" class="space-y-6">
      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">日期范围</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-stone-600 mb-1">开始日期</label>
            <input v-model="startDate" type="date" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
          <div>
            <label class="block text-sm text-stone-600 mb-1">结束日期</label>
            <input v-model="endDate" type="date" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">餐次选择</h2>
        <div class="flex gap-3">
          <button
            v-for="type in [{ key: 'breakfast', label: '早餐' }, { key: 'lunch', label: '午餐' }, { key: 'dinner', label: '晚餐' }]"
            :key="type.key"
            type="button"
            @click="toggleMealType(type.key)"
            :class="[
              'px-4 py-2 rounded-lg text-sm transition-colors',
              mealTypes.includes(type.key) ? 'bg-orange-600 text-white' : 'bg-stone-100 text-stone-600 hover:bg-stone-200'
            ]"
          >{{ type.label }}</button>
        </div>
      </div>

      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">每餐菜品数量</h2>
        <div v-if="mealTypes.includes('lunch')" class="space-y-2">
          <h3 class="text-sm text-stone-600">午餐</h3>
          <div class="grid grid-cols-3 gap-3">
            <div>
              <label class="block text-xs text-stone-500">荤菜</label>
              <input v-model.number="lunchMeat" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
            <div>
              <label class="block text-xs text-stone-500">素菜</label>
              <input v-model.number="lunchVeg" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
            <div>
              <label class="block text-xs text-stone-500">汤</label>
              <input v-model.number="lunchSoup" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
          </div>
        </div>
        <div v-if="mealTypes.includes('dinner')" class="space-y-2">
          <h3 class="text-sm text-stone-600">晚餐</h3>
          <div class="grid grid-cols-3 gap-3">
            <div>
              <label class="block text-xs text-stone-500">荤菜</label>
              <input v-model.number="dinnerMeat" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
            <div>
              <label class="block text-xs text-stone-500">素菜</label>
              <input v-model.number="dinnerVeg" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
            <div>
              <label class="block text-xs text-stone-500">汤</label>
              <input v-model.number="dinnerSoup" type="number" min="0" max="5" class="w-full px-2 py-1.5 border border-stone-300 rounded-lg text-sm" />
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">偏好设置</h2>
        <div>
          <label class="block text-sm text-stone-600 mb-1">口味偏好</label>
          <select v-model="tastePreference" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500">
            <option value="">不限</option>
            <option value="清淡">清淡</option>
            <option value="咸">咸</option>
            <option value="辣">辣</option>
          </select>
        </div>
        <div>
          <label class="block text-sm text-stone-600 mb-1">优先使用食材（逗号分隔）</label>
          <input v-model="preferIngredients" type="text" placeholder="如：菠菜,豆腐" class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
        </div>
        <div>
          <label class="block text-sm text-stone-600 mb-1">排除食材（逗号分隔）</label>
          <input v-model="excludeIngredients" type="text" placeholder="如：花生,虾" class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
        </div>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 font-medium transition-colors"
      >
        {{ loading ? '生成中...' : '智能生成菜单' }}
      </button>
    </form>
  </div>
</template>
