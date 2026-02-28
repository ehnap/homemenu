<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getLLMSettings, updateLLMSettings } from '../api/settings'
import type { LLMSettings } from '../types'

const settings = ref<LLMSettings>({
  base_url: '',
  api_key: '',
  model: '',
})
const loading = ref(false)
const saving = ref(false)
const message = ref('')
const error = ref('')

onMounted(async () => {
  loading.value = true
  try {
    settings.value = await getLLMSettings()
  } catch (e: any) {
    error.value = e.message || '加载设置失败'
  } finally {
    loading.value = false
  }
})

async function save() {
  saving.value = true
  error.value = ''
  message.value = ''
  try {
    await updateLLMSettings(settings.value)
    message.value = '保存成功'
    // Refresh to get masked key
    settings.value = await getLLMSettings()
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold text-stone-800 mb-6">设置</h1>

    <div v-if="loading" class="text-center py-12 text-stone-500">加载中...</div>

    <form v-else @submit.prevent="save" class="space-y-6">
      <div class="bg-white rounded-xl p-6 border border-stone-200 space-y-4">
        <h2 class="font-semibold text-stone-800">LLM 配置</h2>
        <p class="text-sm text-stone-500">配置 LLM API 后可使用智能菜单生成和菜谱文本解析功能。支持 OpenAI 兼容的 API。</p>

        <div>
          <label class="block text-sm text-stone-600 mb-1">API Base URL</label>
          <input
            v-model="settings.base_url"
            type="text"
            placeholder="https://api.openai.com"
            class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
        </div>

        <div>
          <label class="block text-sm text-stone-600 mb-1">API Key</label>
          <input
            v-model="settings.api_key"
            type="password"
            placeholder="sk-..."
            class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
        </div>

        <div>
          <label class="block text-sm text-stone-600 mb-1">模型名称</label>
          <input
            v-model="settings.model"
            type="text"
            placeholder="gpt-4o-mini"
            class="w-full px-3 py-2 border border-stone-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
        </div>
      </div>

      <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg">{{ error }}</div>
      <div v-if="message" class="bg-green-50 text-green-600 text-sm p-3 rounded-lg">{{ message }}</div>

      <button
        type="submit"
        :disabled="saving"
        class="w-full py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50 font-medium transition-colors"
      >
        {{ saving ? '保存中...' : '保存设置' }}
      </button>
    </form>
  </div>
</template>
