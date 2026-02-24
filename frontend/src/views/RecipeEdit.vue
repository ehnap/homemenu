<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRecipeStore } from '../stores/recipe'
import { uploadImage } from '../api/recipe'
import type { Ingredient, Step } from '../types'
import IngredientInput from '../components/IngredientInput.vue'
import ImageUpload from '../components/ImageUpload.vue'

const route = useRoute()
const router = useRouter()
const store = useRecipeStore()

const isEdit = computed(() => !!route.params.id)
const id = computed(() => Number(route.params.id))

const name = ref('')
const cookTime = ref(0)
const difficulty = ref('')
const tags = ref<string[]>([])
const tagInput = ref('')
const coverImage = ref('')
const calories = ref(0)
const notes = ref('')
const ingredients = ref<Ingredient[]>([{ name: '', amount: '', unit: '' }])
const steps = ref<Step[]>([{ order: 1, description: '' }])
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  if (isEdit.value) {
    await store.fetchRecipe(id.value)
    if (store.currentRecipe) {
      name.value = store.currentRecipe.name
      cookTime.value = store.currentRecipe.cook_time
      difficulty.value = store.currentRecipe.difficulty
      tags.value = [...store.currentRecipe.tags]
      coverImage.value = store.currentRecipe.cover_image
      calories.value = store.currentRecipe.calories
      notes.value = store.currentRecipe.notes
      ingredients.value = store.currentRecipe.ingredients.length > 0
        ? store.currentRecipe.ingredients.map(i => ({ ...i }))
        : [{ name: '', amount: '', unit: '' }]
      steps.value = store.currentRecipe.steps.length > 0
        ? store.currentRecipe.steps.map(s => ({ ...s }))
        : [{ order: 1, description: '' }]
    }
  }
})

function addIngredient() {
  ingredients.value.push({ name: '', amount: '', unit: '' })
}

function removeIngredient(index: number) {
  if (ingredients.value.length > 1) {
    ingredients.value.splice(index, 1)
  }
}

function addStep() {
  steps.value.push({ order: steps.value.length + 1, description: '' })
}

function removeStep(index: number) {
  steps.value.splice(index, 1)
  steps.value.forEach((s, i) => (s.order = i + 1))
}

function addTag() {
  const t = tagInput.value.trim()
  if (t && !tags.value.includes(t)) {
    tags.value.push(t)
  }
  tagInput.value = ''
}

function removeTag(index: number) {
  tags.value.splice(index, 1)
}

async function onCoverUpload(file: File) {
  const result = await uploadImage(file)
  coverImage.value = result.url
}

async function submit() {
  if (!name.value.trim()) {
    error.value = '请输入菜名'
    return
  }

  const validIngredients = ingredients.value.filter(i => i.name.trim())
  if (validIngredients.length === 0) {
    error.value = '请至少添加一个食材'
    return
  }

  error.value = ''
  loading.value = true

  const data = {
    name: name.value,
    cook_time: cookTime.value,
    difficulty: difficulty.value,
    tags: tags.value,
    cover_image: coverImage.value,
    calories: calories.value,
    notes: notes.value,
    ingredients: validIngredients,
    steps: steps.value.filter(s => s.description.trim()),
  }

  try {
    if (isEdit.value) {
      await store.editRecipe(id.value, data)
      router.push(`/recipes/${id.value}`)
    } else {
      const recipe = await store.addRecipe(data)
      router.push(`/recipes/${recipe.id}`)
    }
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-stone-800">{{ isEdit ? '编辑菜谱' : '添加菜谱' }}</h1>
      <button @click="router.back()" class="text-sm text-stone-500 hover:text-stone-700">取消</button>
    </div>

    <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg mb-4">{{ error }}</div>

    <form @submit.prevent="submit" class="space-y-6">
      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">基本信息</h2>

        <div>
          <label class="block text-sm text-stone-600 mb-1">菜名 *</label>
          <input v-model="name" type="text" required class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500" />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm text-stone-600 mb-1">烹饪时间（分钟）</label>
            <input v-model.number="cookTime" type="number" min="0" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
          <div>
            <label class="block text-sm text-stone-600 mb-1">难度</label>
            <select v-model="difficulty" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500">
              <option value="">不选</option>
              <option value="简单">简单</option>
              <option value="中等">中等</option>
              <option value="复杂">复杂</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-stone-600 mb-1">卡路里</label>
            <input v-model.number="calories" type="number" min="0" class="w-full px-3 py-2 border border-stone-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
        </div>

        <div>
          <label class="block text-sm text-stone-600 mb-1">标签</label>
          <div class="flex flex-wrap gap-2 mb-2">
            <span v-for="(tag, i) in tags" :key="i" class="px-2 py-1 bg-orange-50 text-orange-600 text-xs rounded-full flex items-center gap-1">
              {{ tag }}
              <button type="button" @click="removeTag(i)" class="hover:text-orange-800">&times;</button>
            </span>
          </div>
          <div class="flex gap-2">
            <input v-model="tagInput" @keydown.enter.prevent="addTag" placeholder="输入标签后回车" class="flex-1 px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
            <button type="button" @click="addTag" class="px-3 py-2 bg-stone-100 text-stone-700 rounded-lg text-sm hover:bg-stone-200">添加</button>
          </div>
        </div>

        <div>
          <label class="block text-sm text-stone-600 mb-1">封面图</label>
          <ImageUpload :current="coverImage" @upload="onCoverUpload" />
        </div>

        <div>
          <label class="block text-sm text-stone-600 mb-1">备注</label>
          <textarea v-model="notes" rows="2" class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"></textarea>
        </div>
      </div>

      <div class="bg-white rounded-xl p-6 border border-stone-200">
        <div class="flex justify-between items-center mb-4">
          <h2 class="font-semibold text-stone-800">食材 *</h2>
          <button type="button" @click="addIngredient" class="text-sm text-orange-600 hover:text-orange-700">+ 添加食材</button>
        </div>
        <div class="space-y-2">
          <IngredientInput
            v-for="(_ing, i) in ingredients"
            :key="i"
            v-model="ingredients[i]!"
            :can-remove="ingredients.length > 1"
            @remove="removeIngredient(i)"
          />
        </div>
      </div>

      <div class="bg-white rounded-xl p-6 border border-stone-200">
        <div class="flex justify-between items-center mb-4">
          <h2 class="font-semibold text-stone-800">做法步骤</h2>
          <button type="button" @click="addStep" class="text-sm text-orange-600 hover:text-orange-700">+ 添加步骤</button>
        </div>
        <div class="space-y-3">
          <div v-for="(step, i) in steps" :key="i" class="flex gap-3 items-start">
            <span class="w-6 h-6 bg-orange-100 text-orange-600 rounded-full flex items-center justify-center text-sm font-medium shrink-0 mt-2">
              {{ step.order }}
            </span>
            <textarea
              v-model="steps[i]!.description"
              rows="2"
              :placeholder="`步骤 ${step.order}`"
              class="flex-1 px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
            ></textarea>
            <button v-if="steps.length > 1" type="button" @click="removeStep(i)" class="text-stone-400 hover:text-red-500 mt-2">
              &times;
            </button>
          </div>
        </div>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 font-medium transition-colors"
      >
        {{ loading ? '保存中...' : '保存菜谱' }}
      </button>
    </form>
  </div>
</template>
