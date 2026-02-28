<script setup lang="ts">
import { ref } from 'vue'
import { parseRecipeText } from '../api/recipe'
import type { Recipe } from '../types'

const emit = defineEmits<{
  parsed: [recipe: Partial<Recipe>]
  close: []
}>()

const text = ref('')
const loading = ref(false)
const error = ref('')

async function parse() {
  if (!text.value.trim()) {
    error.value = '请输入菜谱文本'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const recipe = await parseRecipeText(text.value)
    emit('parsed', recipe)
  } catch (e: any) {
    error.value = e.message || '解析失败，请检查 LLM 配置'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-xl w-full max-w-2xl max-h-[90vh] flex flex-col">
      <div class="flex justify-between items-center p-6 border-b border-stone-200">
        <h2 class="text-lg font-semibold text-stone-800">智能导入</h2>
        <button @click="emit('close')" class="text-stone-400 hover:text-stone-600 text-xl">&times;</button>
      </div>

      <div class="p-6 flex-1 overflow-y-auto space-y-4">
        <p class="text-sm text-stone-500">粘贴菜谱文本，AI 将自动提取菜名、食材、步骤等信息。</p>

        <textarea
          v-model="text"
          rows="12"
          placeholder="例如：&#10;番茄炒蛋&#10;食材：番茄2个、鸡蛋3个、盐适量、糖少许、葱花适量&#10;做法：&#10;1. 番茄切块，鸡蛋打散加少许盐搅匀&#10;2. 热锅凉油，倒入蛋液炒至凝固盛出&#10;3. 锅中加油，放入番茄块翻炒出汁&#10;4. 加入炒好的鸡蛋，加盐和糖调味&#10;5. 撒上葱花即可出锅"
          class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500 resize-none"
          :disabled="loading"
        ></textarea>

        <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg">{{ error }}</div>
      </div>

      <div class="p-6 border-t border-stone-200">
        <button
          @click="parse"
          :disabled="loading || !text.trim()"
          class="w-full py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 font-medium transition-colors"
        >
          {{ loading ? '解析中...' : '智能解析' }}
        </button>
      </div>
    </div>
  </div>
</template>
