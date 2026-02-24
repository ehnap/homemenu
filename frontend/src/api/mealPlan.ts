import client from './client'
import type { MealPlan, MealPlanItem, PlanConfig, ShoppingList } from '../types'

export function listMealPlans(): Promise<MealPlan[]> {
  return client.get<any, MealPlan[]>('/meal-plans')
}

export function getMealPlan(id: number): Promise<MealPlan> {
  return client.get<any, MealPlan>(`/meal-plans/${id}`)
}

export function createMealPlan(plan: Partial<MealPlan>): Promise<MealPlan> {
  return client.post<any, MealPlan>('/meal-plans', plan)
}

export function updateMealPlan(id: number, plan: Partial<MealPlan>): Promise<MealPlan> {
  return client.put<any, MealPlan>(`/meal-plans/${id}`, plan)
}

export function deleteMealPlan(id: number): Promise<void> {
  return client.delete(`/meal-plans/${id}`)
}

export function generateMealPlan(startDate: string, endDate: string, config: PlanConfig): Promise<MealPlan> {
  return client.post<any, MealPlan>('/meal-plans/generate', {
    start_date: startDate,
    end_date: endDate,
    config,
  })
}

export function updateMealPlanItems(planId: number, items: MealPlanItem[]): Promise<MealPlan> {
  return client.put<any, MealPlan>(`/meal-plans/${planId}/items`, { items })
}

export function rerollItem(planId: number, itemId: number): Promise<MealPlanItem> {
  return client.post<any, MealPlanItem>(`/meal-plans/${planId}/items/${itemId}/reroll`)
}

export function getShoppingList(planId: number, mode: 'daily' | 'weekly' = 'weekly'): Promise<ShoppingList> {
  return client.get<any, ShoppingList>(`/meal-plans/${planId}/shopping-list`, { params: { mode } })
}

export function getSharedMealPlan(token: string): Promise<{ meal_plan: MealPlan; shopping_list: ShoppingList }> {
  return client.get<any, { meal_plan: MealPlan; shopping_list: ShoppingList }>(`/share/${token}`)
}
