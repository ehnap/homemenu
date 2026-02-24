import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listRecipes, getRecipe, createRecipe, updateRecipe, deleteRecipe } from '../api/recipe'
import type { Recipe } from '../types'

export const useRecipeStore = defineStore('recipe', () => {
  const recipes = ref<Recipe[]>([])
  const currentRecipe = ref<Recipe | null>(null)
  const loading = ref(false)

  async function fetchRecipes(params?: Record<string, string>) {
    loading.value = true
    try {
      recipes.value = await listRecipes(params)
    } finally {
      loading.value = false
    }
  }

  async function fetchRecipe(id: number) {
    loading.value = true
    try {
      currentRecipe.value = await getRecipe(id)
    } finally {
      loading.value = false
    }
  }

  async function addRecipe(recipe: Partial<Recipe>) {
    const created = await createRecipe(recipe)
    recipes.value.unshift(created)
    return created
  }

  async function editRecipe(id: number, recipe: Partial<Recipe>) {
    const updated = await updateRecipe(id, recipe)
    const idx = recipes.value.findIndex((r) => r.id === id)
    if (idx !== -1) recipes.value[idx] = updated
    return updated
  }

  async function removeRecipe(id: number) {
    await deleteRecipe(id)
    recipes.value = recipes.value.filter((r) => r.id !== id)
  }

  return { recipes, currentRecipe, loading, fetchRecipes, fetchRecipe, addRecipe, editRecipe, removeRecipe }
})
