<template>
  <div
    class="meal-slot-card"
    :class="{ 'meal-slot-card--empty': !props.mealSlot.recipe }"
    @click="$emit('click', props.mealSlot)"
  >
    <div class="slot-header">
      <div class="meal-type-icon" :class="{ 'meal-type-icon--empty': !props.mealSlot.recipe }">
        <q-icon
          :name="getMealIcon(props.mealSlot.mealType)"
          size="14px"
          :color="props.mealSlot.recipe ? 'white' : 'grey-5'"
        />
      </div>
      <span class="meal-type-label">{{ props.mealSlot.mealType }}</span>
      <q-btn
        v-if="props.mealSlot.recipe"
        flat
        dense
        round
        size="xs"
        icon="close"
        class="clear-btn"
        @click.stop="$emit('clear', props.mealSlot)"
      />
    </div>

    <template v-if="props.mealSlot.recipe">
      <p class="recipe-name">{{ props.mealSlot.recipe.name }}</p>
    </template>
    <template v-else>
      <p class="empty-text">
        <q-icon name="add" size="14px" />
        Add recipe
      </p>
    </template>
  </div>
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
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&display=swap');

.meal-slot-card {
  background: white;
  border-radius: 14px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid rgba(45, 31, 26, 0.06);
  box-shadow: 0 2px 8px rgba(45, 31, 26, 0.04);
  min-height: 72px;
  position: relative;

  &:active {
    transform: scale(0.98);
  }

  @media (hover: hover) {
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(255, 127, 80, 0.12);
    }
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

.slot-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
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
    background: #e8e2df;
  }
}

.meal-type-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 600;
  color: #a8a0a0;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  flex: 1;
}

.clear-btn {
  color: #a8a0a0;

  &:hover {
    color: #ff6347;
  }
}

.recipe-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-text {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 500;
  color: #a8a0a0;
  margin: 0;
}
</style>
