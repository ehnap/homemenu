import client from './client'
import type { LLMSettings } from '../types'

export function getLLMSettings(): Promise<LLMSettings> {
  return client.get<any, LLMSettings>('/settings/llm')
}

export function updateLLMSettings(settings: LLMSettings): Promise<void> {
  return client.put('/settings/llm', settings)
}
