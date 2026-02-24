import client from './client'
import type { TokenPair } from '../types'

export function register(username: string, password: string, nickname: string) {
  return client.post<any, any>('/auth/register', { username, password, nickname })
}

export function login(username: string, password: string): Promise<TokenPair> {
  return client.post<any, TokenPair>('/auth/login', { username, password })
}

export function refreshToken(refresh_token: string): Promise<TokenPair> {
  return client.post<any, TokenPair>('/auth/refresh', { refresh_token })
}
