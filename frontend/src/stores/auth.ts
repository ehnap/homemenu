import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const tokens = await apiLogin(username, password)
    token.value = tokens.access_token
    localStorage.setItem('access_token', tokens.access_token)
    localStorage.setItem('refresh_token', tokens.refresh_token)
  }

  async function register(username: string, password: string, nickname: string) {
    await apiRegister(username, password, nickname)
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  return { token, isLoggedIn, login, register, logout }
})
