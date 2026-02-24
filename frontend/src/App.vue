<script setup lang="ts">
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-stone-50">
    <nav v-if="auth.isLoggedIn" class="bg-white shadow-sm border-b border-stone-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-14 items-center">
          <div class="flex items-center space-x-8">
            <router-link to="/" class="text-lg font-bold text-orange-600">家味</router-link>
            <router-link
              to="/recipes"
              class="text-sm text-stone-600 hover:text-orange-600 transition-colors"
              active-class="text-orange-600 font-medium"
            >菜谱</router-link>
            <router-link
              to="/meal-plans"
              class="text-sm text-stone-600 hover:text-orange-600 transition-colors"
              active-class="text-orange-600 font-medium"
            >菜单</router-link>
          </div>
          <button
            @click="logout"
            class="text-sm text-stone-500 hover:text-stone-700 transition-colors"
          >退出</button>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <router-view />
    </main>
  </div>
</template>
