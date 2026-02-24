<script setup lang="ts">
const TAGS = ['咸', '辣', '清淡', '甜', '酸', '荤', '素', '汤']

defineProps<{
  selected: string
}>()

const emit = defineEmits<{
  change: [tag: string]
}>()
</script>

<template>
  <div class="flex gap-1.5 flex-wrap">
    <button
      @click="emit('change', '')"
      :class="['px-2.5 py-1 rounded-full text-xs transition-colors',
        !selected ? 'bg-orange-600 text-white' : 'bg-stone-100 text-stone-600 hover:bg-stone-200']"
    >全部</button>
    <button
      v-for="tag in TAGS"
      :key="tag"
      @click="emit('change', tag)"
      :class="['px-2.5 py-1 rounded-full text-xs transition-colors',
        selected === tag ? 'bg-orange-600 text-white' : 'bg-stone-100 text-stone-600 hover:bg-stone-200']"
    >{{ tag }}</button>
  </div>
</template>
