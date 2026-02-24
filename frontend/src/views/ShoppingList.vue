<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getShoppingList } from '../api/mealPlan'
import type { ShoppingList } from '../types'

const route = useRoute()
const router = useRouter()
const id = Number(route.params.id)
const mode = ref<'daily' | 'weekly'>('weekly')
const list = ref<ShoppingList | null>(null)
const loading = ref(false)
const checkedItems = ref<Set<string>>(new Set())

onMounted(() => {
  loadCheckedItems()
  fetchList()
})

async function fetchList() {
  loading.value = true
  try {
    list.value = await getShoppingList(id, mode.value)
  } finally {
    loading.value = false
  }
}

function switchMode(m: 'daily' | 'weekly') {
  mode.value = m
  fetchList()
}

function toggleChecked(key: string) {
  if (checkedItems.value.has(key)) {
    checkedItems.value.delete(key)
  } else {
    checkedItems.value.add(key)
  }
  saveCheckedItems()
}

function saveCheckedItems() {
  localStorage.setItem(`shopping_checked_${id}`, JSON.stringify([...checkedItems.value]))
}

function loadCheckedItems() {
  const saved = localStorage.getItem(`shopping_checked_${id}`)
  if (saved) {
    checkedItems.value = new Set(JSON.parse(saved))
  }
}

function formatDate(date: string) {
  const d = new Date(date)
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return `${date.substring(5).replace('-', '/')} ${weekdays[d.getDay()]}`
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <div>
        <button @click="router.back()" class="text-sm text-stone-500 hover:text-stone-700 mb-1">&larr; 返回</button>
        <h1 class="text-2xl font-bold text-stone-800">购物清单</h1>
      </div>
      <div class="flex gap-1 bg-stone-100 rounded-lg p-0.5">
        <button
          @click="switchMode('weekly')"
          :class="['px-3 py-1.5 rounded-md text-sm transition-colors', mode === 'weekly' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500']"
        >一周汇总</button>
        <button
          @click="switchMode('daily')"
          :class="['px-3 py-1.5 rounded-md text-sm transition-colors', mode === 'daily' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500']"
        >按天查看</button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-stone-500">加载中...</div>

    <template v-else-if="list">
      <!-- Weekly view -->
      <div v-if="mode === 'weekly' && list.weekly" class="bg-white rounded-xl border border-stone-200">
        <div v-for="(item, i) in list.weekly" :key="i"
          :class="['flex items-center justify-between px-4 py-3', i < list.weekly.length - 1 ? 'border-b border-stone-100' : '']"
        >
          <div class="flex items-center gap-3">
            <button
              @click="toggleChecked(`w_${item.name}_${item.unit}`)"
              :class="['w-5 h-5 rounded border-2 flex items-center justify-center transition-colors',
                checkedItems.has(`w_${item.name}_${item.unit}`) ? 'bg-green-500 border-green-500 text-white' : 'border-stone-300']"
            >
              <span v-if="checkedItems.has(`w_${item.name}_${item.unit}`)" class="text-xs">&#10003;</span>
            </button>
            <span :class="['text-sm', checkedItems.has(`w_${item.name}_${item.unit}`) ? 'text-stone-400 line-through' : 'text-stone-700']">
              {{ item.name }}
            </span>
          </div>
          <span class="text-sm text-stone-500">{{ item.amount }}{{ item.unit }}</span>
        </div>
        <div v-if="list.weekly.length === 0" class="text-center py-8 text-stone-500 text-sm">暂无食材</div>
      </div>

      <!-- Daily view -->
      <div v-if="mode === 'daily' && list.daily" class="space-y-4">
        <div v-for="day in list.daily" :key="day.date" class="bg-white rounded-xl border border-stone-200">
          <div class="px-4 py-3 border-b border-stone-100">
            <h3 class="font-medium text-stone-800 text-sm">{{ formatDate(day.date) }}</h3>
          </div>
          <div v-for="(item, i) in day.items" :key="i"
            :class="['flex items-center justify-between px-4 py-2.5', i < day.items.length - 1 ? 'border-b border-stone-50' : '']"
          >
            <div class="flex items-center gap-3">
              <button
                @click="toggleChecked(`d_${day.date}_${item.name}_${item.unit}`)"
                :class="['w-5 h-5 rounded border-2 flex items-center justify-center transition-colors',
                  checkedItems.has(`d_${day.date}_${item.name}_${item.unit}`) ? 'bg-green-500 border-green-500 text-white' : 'border-stone-300']"
              >
                <span v-if="checkedItems.has(`d_${day.date}_${item.name}_${item.unit}`)" class="text-xs">&#10003;</span>
              </button>
              <span :class="['text-sm', checkedItems.has(`d_${day.date}_${item.name}_${item.unit}`) ? 'text-stone-400 line-through' : 'text-stone-700']">
                {{ item.name }}
              </span>
            </div>
            <span class="text-sm text-stone-500">{{ item.amount }}{{ item.unit }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
