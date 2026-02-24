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
