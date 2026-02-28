<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRecipeStore } from '../stores/recipe'
import { generateShareToken } from '../api/recipe'

const route = useRoute()
const router = useRouter()
const store = useRecipeStore()

const id = Number(route.params.id)
const shareLink = ref('')
const showShareModal = ref(false)
const shareLoading = ref(false)
const copySuccess = ref(false)

onMounted(() => {
  store.fetchRecipe(id)
})

async function handleDelete() {
  if (!confirm('确定要删除这个菜谱吗？')) return
  await store.removeRecipe(id)
  router.push('/recipes')
}

async function handleShare() {
  shareLoading.value = true
  try {
    const { share_token } = await generateShareToken(id)
    shareLink.value = `${window.location.origin}/share/recipe/${share_token}`
    showShareModal.value = true
  } catch (e) {
    alert('生成分享链接失败')
  } finally {
    shareLoading.value = false
  }
}

async function copyShareLink() {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(shareLink.value)
    } else {
      // Fallback for non-HTTPS (e.g. LAN IP access)
      const textarea = document.createElement('textarea')
      textarea.value = shareLink.value
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    copySuccess.value = true
    setTimeout(() => { copySuccess.value = false }, 2000)
  } catch {
    alert('复制失败，请手动选中链接复制')
  }
}
</script>

<template>
  <div v-if="store.loading" class="text-center py-12 text-stone-500">加载中...</div>

  <div v-else-if="store.currentRecipe" class="max-w-3xl mx-auto">
    <div class="flex justify-between items-start mb-6">
      <div>
        <button @click="router.back()" class="text-sm text-stone-500 hover:text-stone-700 mb-2">&larr; 返回</button>
        <h1 class="text-2xl font-bold text-stone-800">{{ store.currentRecipe.name }}</h1>
        <div class="flex gap-2 mt-2">
          <span v-if="store.currentRecipe.difficulty" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
            {{ store.currentRecipe.difficulty }}
          </span>
          <span v-if="store.currentRecipe.cook_time" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
            {{ store.currentRecipe.cook_time }}分钟
          </span>
          <span v-if="store.currentRecipe.calories" class="px-2 py-0.5 bg-stone-100 text-stone-600 text-xs rounded-full">
            {{ store.currentRecipe.calories }}卡
          </span>
          <span v-for="tag in store.currentRecipe.tags" :key="tag" class="px-2 py-0.5 bg-orange-50 text-orange-600 text-xs rounded-full">
            {{ tag }}
          </span>
        </div>
      </div>
      <div class="flex gap-2">
        <button
          @click="handleShare"
          :disabled="shareLoading"
          class="px-3 py-1.5 bg-orange-50 text-orange-600 rounded-lg text-sm hover:bg-orange-100 disabled:opacity-50"
        >{{ shareLoading ? '生成中...' : '分享' }}</button>
        <router-link
          :to="`/recipes/${id}/edit`"
          class="px-3 py-1.5 bg-stone-100 text-stone-700 rounded-lg text-sm hover:bg-stone-200"
        >编辑</router-link>
        <button
          @click="handleDelete"
          class="px-3 py-1.5 bg-red-50 text-red-600 rounded-lg text-sm hover:bg-red-100"
        >删除</button>
      </div>
    </div>

    <img
      v-if="store.currentRecipe.cover_image"
      :src="store.currentRecipe.cover_image"
      :alt="store.currentRecipe.name"
      class="w-full h-64 object-cover rounded-xl mb-6"
    />

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="md:col-span-1">
        <div class="bg-white rounded-xl p-4 border border-stone-200">
          <h2 class="font-semibold text-stone-800 mb-3">食材</h2>
          <ul class="space-y-2">
            <li v-for="ing in store.currentRecipe.ingredients" :key="ing.id" class="flex justify-between text-sm">
              <span class="text-stone-700">{{ ing.name }}</span>
              <span class="text-stone-500">{{ ing.amount }}{{ ing.unit }}</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="md:col-span-2">
        <div class="bg-white rounded-xl p-4 border border-stone-200">
          <h2 class="font-semibold text-stone-800 mb-3">做法步骤</h2>
          <div v-if="store.currentRecipe.steps.length === 0" class="text-stone-500 text-sm">暂无步骤</div>
          <ol class="space-y-4">
            <li v-for="step in store.currentRecipe.steps" :key="step.order" class="flex gap-3">
              <span class="w-6 h-6 bg-orange-100 text-orange-600 rounded-full flex items-center justify-center text-sm font-medium shrink-0">
                {{ step.order }}
              </span>
              <div>
                <p class="text-sm text-stone-700">{{ step.description }}</p>
                <img v-if="step.image_url" :src="step.image_url" class="mt-2 rounded-lg max-h-40" />
              </div>
            </li>
          </ol>
        </div>

        <div v-if="store.currentRecipe.notes" class="bg-white rounded-xl p-4 border border-stone-200 mt-4">
          <h2 class="font-semibold text-stone-800 mb-2">备注</h2>
          <p class="text-sm text-stone-600">{{ store.currentRecipe.notes }}</p>
        </div>
      </div>
    </div>
  </div>

  <!-- Share modal -->
  <div v-if="showShareModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showShareModal = false">
    <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-stone-800">分享菜谱</h3>
        <button @click="showShareModal = false" class="text-stone-400 hover:text-stone-600 text-xl">&times;</button>
      </div>
      <p class="text-sm text-stone-600 mb-3">复制以下链接分享给他人，无需登录即可查看：</p>
      <div class="flex gap-2">
        <input
          :value="shareLink"
          readonly
          class="flex-1 px-3 py-2 border border-stone-300 rounded-lg text-sm bg-stone-50"
          @focus="($event.target as HTMLInputElement).select()"
        />
        <button
          @click="copyShareLink"
          :class="['px-4 py-2 rounded-lg text-sm', copySuccess ? 'bg-green-500 text-white' : 'bg-orange-500 text-white hover:bg-orange-600']"
        >{{ copySuccess ? '已复制' : '复制' }}</button>
      </div>
    </div>
  </div>
</template>
