<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const nickname = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (isRegister.value) {
      await auth.register(username.value, password.value, nickname.value)
      await auth.login(username.value, password.value)
    } else {
      await auth.login(username.value, password.value)
    }
    router.push('/recipes')
  } catch (e: any) {
    error.value = e.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-stone-50">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-orange-600">家味</h1>
        <p class="text-stone-500 mt-1">家庭菜谱与菜单规划</p>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-stone-200">
        <h2 class="text-lg font-semibold text-stone-800 mb-4">
          {{ isRegister ? '注册' : '登录' }}
        </h2>

        <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg mb-4">
          {{ error }}
        </div>

        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="block text-sm text-stone-600 mb-1">用户名</label>
            <input
              v-model="username"
              type="text"
              required
              class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm text-stone-600 mb-1">密码</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-transparent"
            />
          </div>

          <div v-if="isRegister">
            <label class="block text-sm text-stone-600 mb-1">昵称</label>
            <input
              v-model="nickname"
              type="text"
              class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-transparent"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 transition-colors"
          >
            {{ loading ? '请稍候...' : (isRegister ? '注册' : '登录') }}
          </button>
        </form>

        <p class="mt-4 text-center text-sm text-stone-500">
          <button @click="isRegister = !isRegister" class="text-orange-600 hover:underline">
            {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
          </button>
        </p>
      </div>
    </div>
  </div>
</template>
