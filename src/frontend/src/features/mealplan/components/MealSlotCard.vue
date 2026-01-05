<template>
  <q-card
    class="meal-slot-card"
    :class="{ 'meal-slot-card--empty': !props.mealSlot.recipe }"
    flat
    @click="$emit('click', props.mealSlot)"
  >
    <q-card-section class="q-pa-sm">
      <div class="tw-flex tw-items-center tw-gap-2 tw-mb-1">
        <div class="meal-type-icon" :class="{ 'meal-type-icon--empty': !props.mealSlot.recipe }">
          <q-icon :name="getMealIcon(props.mealSlot.mealType)" size="14px" :color="props.mealSlot.recipe ? 'white' : 'grey-5'" />
        </div>
        <span class="text-caption text-weight-medium text-uppercase text-grey-7">
          {{ props.mealSlot.mealType }}
        </span>
      </div>

      <template v-if="props.mealSlot.recipe">
        <div class="text-body2 ellipsis tw-font-medium">
          {{ props.mealSlot.recipe.name }}
        </div>
        <q-btn
          flat
          dense
          round
          size="sm"
          icon="close"
          class="clear-btn"
          @click.stop="$emit('clear', props.mealSlot)"
        />
      </template>

      <template v-else>
        <div class="text-body2 text-grey-5 tw-flex tw-items-center tw-gap-1">
          <q-icon name="add" size="16px" /> Add recipe
        </div>
      </template>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import type { MealSlot, MealType } from '@features/mealplan/types/mealplan';

const props = defineProps<{
  mealSlot: MealSlot;
}>();

defineEmits<{
  click: [slot: MealSlot];
  clear: [slot: MealSlot];
}>();

function getMealIcon(mealType: MealType): string {
  const icons: Record<MealType, string> = {
    breakfast: 'free_breakfast',
    lunch: 'lunch_dining',
    dinner: 'dinner_dining',
    snack: 'cookie',
  };
  return icons[mealType] || 'restaurant';
}
</script>

<style scoped lang="scss">
.meal-slot-card {
  cursor: pointer;
  min-height: 70px;
  position: relative;
  transition: all 0.2s ease;
  border-radius: 12px;
  background: white;
  border: 1px solid rgba(0, 0, 0, 0.06);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 127, 80, 0.15);
  }

  &--empty {
    border: 2px dashed rgba(255, 127, 80, 0.3);
    background: #fffaf8;
    box-shadow: none;

    &:hover {
      border-color: #ff7f50;
      background: #fff5f2;
    }
  }
}

.meal-type-icon {
  width: 24px;
  height: 24px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;

  &--empty {
    background: #f0f0f0;
  }
}

.clear-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  color: #999;

  &:hover {
    color: #ff6347;
  }
}
</style>
