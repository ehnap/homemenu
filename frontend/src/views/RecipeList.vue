<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRecipeStore } from '../stores/recipe'
import RecipeCard from '../components/RecipeCard.vue'
import TagFilter from '../components/TagFilter.vue'

const store = useRecipeStore()
const searchQuery = ref('')
const selectedTag = ref('')
const ingredientSearch = ref('')

onMounted(() => {
  store.fetchRecipes()
})

function applyFilters() {
  const params: Record<string, string> = {}
  if (searchQuery.value) params.q = searchQuery.value
  if (selectedTag.value) params.tag = selectedTag.value
  if (ingredientSearch.value) params.ingredients = ingredientSearch.value
  store.fetchRecipes(params)
}

function onTagChange(tag: string) {
  selectedTag.value = tag
  applyFilters()
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-stone-800">菜谱库</h1>
      <router-link
        to="/recipes/new"
        class="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 text-sm transition-colors"
      >添加菜谱</router-link>
    </div>

    <div class="flex flex-wrap gap-3 mb-6">
      <input
        v-model="searchQuery"
        @input="applyFilters"
        type="text"
        placeholder="搜索菜名..."
        class="px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
      />
      <input
        v-model="ingredientSearch"
        @input="applyFilters"
        type="text"
        placeholder="按食材搜索（逗号分隔）..."
        class="px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
      />
      <TagFilter :selected="selectedTag" @change="onTagChange" />
    </div>

    <div v-if="store.loading" class="text-center py-12 text-stone-500">加载中...</div>

    <div v-else-if="store.recipes.length === 0" class="text-center py-12">
      <p class="text-stone-500">还没有菜谱</p>
      <router-link to="/recipes/new" class="text-orange-600 hover:underline text-sm mt-2 inline-block">
        添加第一个菜谱
      </router-link>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <RecipeCard v-for="recipe in store.recipes" :key="recipe.id" :recipe="recipe" />
    </div>
  </div>
</template>
