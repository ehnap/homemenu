<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMealPlanStore } from '../stores/mealPlan'
import { useRecipeStore } from '../stores/recipe'
import type { MealPlanItem } from '../types'
import { MEAL_TYPE_LABELS, WEEKDAY_LABELS } from '../types'
import MealSlot from '../components/MealSlot.vue'

const route = useRoute()
const router = useRouter()
const planStore = useMealPlanStore()
const recipeStore = useRecipeStore()

const id = Number(route.params.id)
const showRecipePicker = ref(false)
const pickerTarget = ref<{ date: string; mealType: string } | null>(null)
const saving = ref(false)

onMounted(async () => {
  await planStore.fetchPlan(id)
  await recipeStore.fetchRecipes()
})

const plan = computed(() => planStore.currentPlan)

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

async function handleReroll(itemId: number) {
  await planStore.reroll(id, itemId)
}

function removeItem(itemId: number) {
  if (!plan.value?.items) return
  plan.value.items = plan.value.items.filter(i => i.id !== itemId)
}

function openRecipePicker(date: string, mealType: string) {
  pickerTarget.value = { date, mealType }
  showRecipePicker.value = true
}

function addRecipeToPlan(recipeId: number) {
  if (!plan.value?.items || !pickerTarget.value) return
  const items = getItems(pickerTarget.value.date, pickerTarget.value.mealType)
  const newItem: MealPlanItem = {
    id: Date.now(),
    meal_plan_id: id,
    recipe_id: recipeId,
    date: pickerTarget.value.date,
    meal_type: pickerTarget.value.mealType,
    sort_order: items.length,
    recipe: recipeStore.recipes.find(r => r.id === recipeId),
  }
  plan.value.items.push(newItem)
  showRecipePicker.value = false
}

async function saveChanges() {
  if (!plan.value?.items) return
  saving.value = true
  try {
    const items = plan.value.items.map((item, i) => ({
      ...item,
      sort_order: i,
    }))
    await planStore.saveItems(id, items)
  } finally {
    saving.value = false
  }
}

function handleDragChange(date: string, mealType: string, newItems: MealPlanItem[]) {
  if (!plan.value?.items) return
  // Remove old items for this slot
  plan.value.items = plan.value.items.filter(
    i => !(i.date === date && i.meal_type === mealType)
  )
  // Add new items with updated date/meal_type
  newItems.forEach((item, idx) => {
    item.date = date
    item.meal_type = mealType
    item.sort_order = idx
    plan.value!.items!.push(item)
  })
}
</script>

<template>
  <div v-if="planStore.loading" class="text-center py-12 text-stone-500">加载中...</div>

  <div v-else-if="plan">
    <div class="flex justify-between items-center mb-6">
      <div>
        <button @click="router.back()" class="text-sm text-stone-500 hover:text-stone-700 mb-1">&larr; 返回</button>
        <h1 class="text-2xl font-bold text-stone-800">{{ plan.name || '菜单编辑' }}</h1>
      </div>
      <div class="flex gap-2">
        <router-link
          :to="`/meal-plans/${id}/shopping`"
          class="px-4 py-2 bg-green-50 text-green-600 rounded-lg text-sm hover:bg-green-100"
        >购物清单</router-link>
        <button
          @click="saveChanges"
          :disabled="saving"
          class="px-4 py-2 bg-orange-600 text-white rounded-lg text-sm hover:bg-orange-700 disabled:opacity-50"
        >{{ saving ? '保存中...' : '保存修改' }}</button>
      </div>
    </div>

    <div v-if="plan.share_token" class="mb-4 text-sm text-stone-500">
      分享链接：<span class="text-orange-600 select-all">/share/{{ plan.share_token }}</span>
    </div>

    <div class="overflow-x-auto">
      <div class="inline-grid gap-2" :style="{ gridTemplateColumns: `80px repeat(${dates.length}, minmax(140px, 1fr))` }">
        <!-- Header row -->
        <div></div>
        <div v-for="date in dates" :key="date" class="text-center py-2">
          <div class="font-medium text-stone-800 text-sm">{{ getWeekday(date) }}</div>
          <div class="text-xs text-stone-500">{{ formatDate(date) }}</div>
        </div>

        <!-- Meal type rows -->
        <template v-for="mealType in mealTypes" :key="mealType">
          <div class="flex items-center justify-center text-sm font-medium text-stone-600 py-2">
            {{ MEAL_TYPE_LABELS[mealType] || mealType }}
          </div>
          <MealSlot
            v-for="date in dates"
            :key="`${date}-${mealType}`"
            :items="getItems(date, mealType)"
            :date="date"
            :meal-type="mealType"
            @reroll="handleReroll"
            @remove="removeItem"
            @add="openRecipePicker(date, mealType)"
            @update="handleDragChange(date, mealType, $event)"
          />
        </template>
      </div>
    </div>

    <!-- Recipe picker modal -->
    <div v-if="showRecipePicker" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showRecipePicker = false">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4 max-h-[70vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-stone-800">选择菜品</h3>
          <button @click="showRecipePicker = false" class="text-stone-400 hover:text-stone-600">&times;</button>
        </div>
        <div v-if="recipeStore.recipes.length === 0" class="text-stone-500 text-sm text-center py-4">
          暂无菜谱
        </div>
        <div v-else class="space-y-2">
          <button
            v-for="recipe in recipeStore.recipes"
            :key="recipe.id"
            @click="addRecipeToPlan(recipe.id)"
            class="w-full text-left px-3 py-2 hover:bg-orange-50 rounded-lg text-sm text-stone-700 transition-colors"
          >
            {{ recipe.name }}
            <span v-if="recipe.difficulty" class="text-xs text-stone-400 ml-2">{{ recipe.difficulty }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
