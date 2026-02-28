import axios from 'axios'
import client from './client'
import type { Recipe } from '../types'

export function listRecipes(params?: Record<string, string>): Promise<Recipe[]> {
  return client.get<any, Recipe[]>('/recipes', { params })
}

export function getRecipe(id: number): Promise<Recipe> {
  return client.get<any, Recipe>(`/recipes/${id}`)
}

export function createRecipe(recipe: Partial<Recipe>): Promise<Recipe> {
  return client.post<any, Recipe>('/recipes', recipe)
}

export function updateRecipe(id: number, recipe: Partial<Recipe>): Promise<Recipe> {
  return client.put<any, Recipe>(`/recipes/${id}`, recipe)
}

export function deleteRecipe(id: number): Promise<void> {
  return client.delete(`/recipes/${id}`)
}

export function uploadImage(file: File): Promise<{ url: string }> {
  const formData = new FormData()
  formData.append('file', file)
  return client.post<any, { url: string }>('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function parseRecipeText(text: string): Promise<Partial<Recipe>> {
  return client.post<any, Partial<Recipe>>('/recipes/parse-text', { text }, {
    timeout: 60000,
  })
}

export function suggestIngredients(query: string): Promise<string[]> {
  return client.get<any, string[]>('/ingredients/suggestions', {
    params: { q: query },
  })
}

export function generateShareToken(id: number): Promise<{ share_token: string }> {
  return client.post<any, { share_token: string }>(`/recipes/${id}/share`)
}

export function getSharedRecipe(token: string): Promise<Recipe> {
  return axios.get(`/api/share/recipe/${token}`).then(res => {
    const data = res.data
    if (data.code !== 0) {
      throw new Error(data.message || 'Request failed')
    }
    return data.data
  })
}
