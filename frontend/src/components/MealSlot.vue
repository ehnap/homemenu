<script setup lang="ts">
import type { MealPlanItem } from '../types'
import DraggableMealCard from './DraggableMealCard.vue'
import draggable from 'vuedraggable'
import { ref, watch } from 'vue'

const props = defineProps<{
  items: MealPlanItem[]
  date: string
  mealType: string
}>()

const emit = defineEmits<{
  reroll: [itemId: number]
  remove: [itemId: number]
  add: []
  update: [items: MealPlanItem[]]
}>()

const localItems = ref<MealPlanItem[]>([...props.items])

watch(() => props.items, (newItems) => {
  localItems.value = [...newItems]
}, { deep: true })

function onDragChange() {
  emit('update', localItems.value)
}
</script>

<template>
  <div class="bg-white rounded-lg border border-stone-200 p-2 min-h-[80px]">
    <draggable
      v-model="localItems"
      group="meals"
      item-key="id"
      :animation="200"
      class="space-y-1 min-h-[40px]"
      @change="onDragChange"
    >
      <template #item="{ element }">
        <DraggableMealCard
          :item="element"
          @reroll="emit('reroll', element.id)"
          @remove="emit('remove', element.id)"
        />
      </template>
    </draggable>
    <button
      @click="emit('add')"
      class="w-full mt-1 py-1 text-xs text-stone-400 hover:text-orange-600 hover:bg-orange-50 rounded transition-colors"
    >+ 添加</button>
  </div>
</template>
