import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  listMealPlans,
  getMealPlan,
  deleteMealPlan,
  generateMealPlan,
  updateMealPlanItems,
  rerollItem,
} from '../api/mealPlan'
import type { MealPlan, MealPlanItem, PlanConfig } from '../types'

export const useMealPlanStore = defineStore('mealPlan', () => {
  const plans = ref<MealPlan[]>([])
  const currentPlan = ref<MealPlan | null>(null)
  const loading = ref(false)

  async function fetchPlans() {
    loading.value = true
    try {
      plans.value = await listMealPlans()
    } finally {
      loading.value = false
    }
  }

  async function fetchPlan(id: number) {
    loading.value = true
    try {
      currentPlan.value = await getMealPlan(id)
    } finally {
      loading.value = false
    }
  }

  async function removePlan(id: number) {
    await deleteMealPlan(id)
    plans.value = plans.value.filter((p) => p.id !== id)
  }

  async function generate(startDate: string, endDate: string, config: PlanConfig) {
    loading.value = true
    try {
      const plan = await generateMealPlan(startDate, endDate, config)
      plans.value.unshift(plan)
      currentPlan.value = plan
      return plan
    } finally {
      loading.value = false
    }
  }

  async function saveItems(planId: number, items: MealPlanItem[]) {
    const plan = await updateMealPlanItems(planId, items)
    currentPlan.value = plan
    return plan
  }

  async function reroll(planId: number, itemId: number) {
    const item = await rerollItem(planId, itemId)
    if (currentPlan.value?.items) {
      const idx = currentPlan.value.items.findIndex((i) => i.id === itemId)
      if (idx !== -1) {
        currentPlan.value.items[idx] = item
      }
    }
    return item
  }

  return { plans, currentPlan, loading, fetchPlans, fetchPlan, removePlan, generate, saveItems, reroll }
})
