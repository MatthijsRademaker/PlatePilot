<template>
  <q-card
    class="meal-slot-card"
    :class="{ 'meal-slot-card--empty': !mealSlot.recipe }"
    @click="$emit('click', mealSlot)"
  >
    <q-card-section class="q-pa-sm">
      <div class="text-caption text-weight-medium text-uppercase">
        {{ mealSlot.mealType }}
      </div>

      <template v-if="mealSlot.recipe">
        <div class="text-body2 ellipsis q-mt-xs">
          {{ mealSlot.recipe.name }}
        </div>
        <q-btn
          flat
          dense
          round
          size="sm"
          icon="close"
          class="absolute-top-right"
          @click.stop="$emit('clear', mealSlot)"
        />
      </template>

      <template v-else>
        <div class="text-body2 text-grey q-mt-xs">
          <q-icon name="add" /> Add recipe
        </div>
        <q-btn
          flat
          dense
          round
          size="sm"
          icon="auto_awesome"
          color="primary"
          class="absolute-top-right suggest-btn"
          @click.stop="$emit('suggest', mealSlot)"
        >
          <q-tooltip>Get AI suggestions</q-tooltip>
        </q-btn>
      </template>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import type { MealSlot } from '../types';

defineProps<{
  mealSlot: MealSlot;
}>();

defineEmits<{
  click: [slot: MealSlot];
  clear: [slot: MealSlot];
  suggest: [slot: MealSlot];
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

    .suggest-btn {
      opacity: 0.7;
      transition: opacity 0.2s;
    }

    &:hover .suggest-btn {
      opacity: 1;
    }
  }
}
</style>
