export interface User {
  id: number
  username: string
  nickname: string
  created_at: string
}

export interface Step {
  order: number
  description: string
  image_url?: string
}

export interface Ingredient {
  id?: number
  recipe_id?: number
  name: string
  amount: string
  unit: string
}

export interface Recipe {
  id: number
  user_id: number
  name: string
  steps: Step[]
  cook_time: number
  difficulty: string
  tags: string[]
  cover_image: string
  calories: number
  notes: string
  tips: string
  share_token?: string
  ingredients: Ingredient[]
  seasonings: Ingredient[]
  created_at: string
  updated_at: string
}

export interface DishCount {
  meat: number
  vegetable: number
  soup: number
}

export interface PlanConfig {
  meal_types: string[]
  dishes_per_meal?: Record<string, DishCount>
  taste_preference?: string
  prefer_ingredients?: string[]
  exclude_ingredients?: string[]
  use_ai?: boolean
}

export interface MealPlanItem {
  id: number
  meal_plan_id: number
  recipe_id: number
  date: string
  meal_type: string
  sort_order: number
  recipe?: Recipe
}

export interface MealPlan {
  id: number
  user_id: number
  name: string
  start_date: string
  end_date: string
  config: PlanConfig
  share_token?: string
  items?: MealPlanItem[]
  created_at: string
  updated_at: string
}

export interface ShoppingItem {
  name: string
  amount: string
  unit: string
}

export interface DailyShoppingList {
  date: string
  items: ShoppingItem[]
}

export interface ShoppingList {
  daily?: DailyShoppingList[]
  weekly?: ShoppingItem[]
}

export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

export interface LLMSettings {
  base_url: string
  api_key: string
  model: string
}

export const MEAL_TYPE_LABELS: Record<string, string> = {
  breakfast: '早餐',
  lunch: '午餐',
  dinner: '晚餐',
}

export const DIFFICULTY_LABELS: Record<string, string> = {
  '简单': '简单',
  '中等': '中等',
  '复杂': '复杂',
}

export const WEEKDAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
