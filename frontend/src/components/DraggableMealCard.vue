<script setup lang="ts">
import type { MealPlanItem } from '../types'

defineProps<{
  item: MealPlanItem
}>()

const emit = defineEmits<{
  reroll: []
  remove: []
}>()
</script>

<template>
  <div class="bg-orange-50 rounded px-2 py-1.5 cursor-grab active:cursor-grabbing group">
    <div class="flex items-center justify-between">
      <span class="text-xs text-stone-700 truncate flex-1">{{ item.recipe?.name || '未知' }}</span>
      <div class="hidden group-hover:flex gap-1 ml-1">
        <button @click.stop="emit('reroll')" class="text-stone-400 hover:text-orange-600 text-xs" title="换一道">&#x21bb;</button>
        <button @click.stop="emit('remove')" class="text-stone-400 hover:text-red-500 text-xs" title="删除">&times;</button>
      </div>
    </div>
  </div>
</template>
