<template>
  <q-card
    class="meal-slot-card"
    :class="{ 'meal-slot-card--empty': !slot.recipe }"
    @click="$emit('click', slot)"
  >
    <q-card-section class="q-pa-sm">
      <div class="text-caption text-weight-medium text-uppercase">
        {{ slot.mealType }}
      </div>

      <template v-if="slot.recipe">
        <div class="text-body2 ellipsis q-mt-xs">
          {{ slot.recipe.name }}
        </div>
        <q-btn
          flat
          dense
          round
          size="sm"
          icon="close"
          class="absolute-top-right"
          @click.stop="$emit('clear', slot)"
        />
      </template>

      <template v-else>
        <div class="text-body2 text-grey q-mt-xs">
          <q-icon name="add" /> Add recipe
        </div>
      </template>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import type { MealSlot } from '../types';

defineProps<{
  slot: MealSlot;
}>();

defineEmits<{
  click: [slot: MealSlot];
  clear: [slot: MealSlot];
}>();
</script>

<style scoped lang="scss">
.meal-slot-card {
  cursor: pointer;
  min-height: 70px;
  position: relative;
  transition: all 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  &--empty {
    border: 2px dashed #ccc;
    background: transparent;
    box-shadow: none;
  }
}
</style>
