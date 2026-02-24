<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMealPlanStore } from '../stores/mealPlan'

const store = useMealPlanStore()
const router = useRouter()

onMounted(() => {
  store.fetchPlans()
})

async function handleDelete(id: number) {
  if (!confirm('确定要删除这个菜单吗？')) return
  await store.removePlan(id)
}

function formatDate(date: string) {
  return date.substring(5).replace('-', '/')
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-stone-800">周菜单</h1>
      <router-link
        to="/meal-plans/generate"
        class="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 text-sm transition-colors"
      >生成菜单</router-link>
    </div>

    <div v-if="store.loading" class="text-center py-12 text-stone-500">加载中...</div>

    <div v-else-if="store.plans.length === 0" class="text-center py-12">
      <p class="text-stone-500">还没有菜单</p>
      <router-link to="/meal-plans/generate" class="text-orange-600 hover:underline text-sm mt-2 inline-block">
        生成第一个菜单
      </router-link>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="plan in store.plans"
        :key="plan.id"
        class="bg-white rounded-xl p-4 border border-stone-200 flex justify-between items-center hover:border-orange-200 transition-colors"
      >
        <div class="cursor-pointer flex-1" @click="router.push(`/meal-plans/${plan.id}`)">
          <h3 class="font-medium text-stone-800">{{ plan.name || '未命名菜单' }}</h3>
          <p class="text-sm text-stone-500 mt-1">
            {{ formatDate(plan.start_date) }} - {{ formatDate(plan.end_date) }}
          </p>
        </div>
        <div class="flex gap-2">
          <router-link
            :to="`/meal-plans/${plan.id}/shopping`"
            class="px-3 py-1.5 bg-green-50 text-green-600 rounded-lg text-sm hover:bg-green-100"
          >购物清单</router-link>
          <router-link
            :to="`/meal-plans/${plan.id}`"
            class="px-3 py-1.5 bg-stone-100 text-stone-700 rounded-lg text-sm hover:bg-stone-200"
          >编辑</router-link>
          <button
            @click="handleDelete(plan.id)"
            class="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-sm hover:bg-red-100"
          >删除</button>
        </div>
      </div>
    </div>
  </div>
</template>
