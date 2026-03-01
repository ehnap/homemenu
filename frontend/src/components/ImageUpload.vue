<script setup lang="ts">
import { ref } from 'vue'
import { uploadImage } from '../api/recipe'

defineProps<{
  current?: string
}>()

const emit = defineEmits<{
  uploaded: [url: string]
}>()

const uploading = ref(false)
const error = ref('')

async function handleFile(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return

  const file = input.files[0]!
  uploading.value = true
  error.value = ''

  try {
    const result = await uploadImage(file)
    emit('uploaded', result.url)
  } catch (err: any) {
    error.value = err.message || '图片上传失败'
  } finally {
    uploading.value = false
    input.value = ''
  }
}
</script>

<template>
  <div>
    <div v-if="current" class="mb-2">
      <img :src="current" class="h-32 rounded-lg object-cover" />
    </div>
    <label class="inline-block px-3 py-2 bg-stone-100 text-stone-700 rounded-lg text-sm cursor-pointer hover:bg-stone-200 transition-colors">
      {{ uploading ? '上传中...' : (current ? '更换图片' : '选择图片') }}
      <input type="file" accept="image/*" class="hidden" @change="handleFile" :disabled="uploading" />
    </label>
    <p v-if="error" class="mt-1 text-xs text-red-500">{{ error }}</p>
  </div>
</template>
