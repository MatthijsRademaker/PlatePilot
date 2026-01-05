<template>
  <div class="calorie-tracker">
    <div class="tw-text-lg tw-font-semibold tw-text-gray-800 tw-mb-4">Daily Calorie Tracker</div>

    <div class="tw-flex tw-items-center tw-gap-4">
      <!-- Main Calorie Circle -->
      <div class="main-progress">
        <q-circular-progress
          :value="caloriePercentage"
          size="80px"
          :thickness="0.15"
          color="orange"
          track-color="grey-3"
          class="main-circle"
        />
        <div class="progress-text">
          <div class="tw-text-lg tw-font-bold tw-text-gray-800">{{ currentCalories }}</div>
          <div class="tw-text-xs tw-text-gray-500">/ {{ targetCalories }} cal</div>
        </div>
      </div>

      <!-- Meal Breakdown -->
      <div class="tw-flex tw-gap-3 tw-flex-1 tw-overflow-x-auto">
        <div v-for="item in mealBreakdown" :key="item.name" class="meal-item">
          <q-circular-progress
            :value="item.percentage"
            size="48px"
            :thickness="0.12"
            :color="item.color"
            track-color="grey-3"
          />
          <div class="meal-item-image">
            <img :src="item.image" :alt="item.name" />
          </div>
          <div class="tw-text-xs tw-text-gray-600 tw-text-center tw-mt-1 tw-truncate tw-max-w-12">
            {{ item.name }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

// Mock data - in a real app this would come from a store/API
const targetCalories = 2200;
const currentCalories = 1480;

const caloriePercentage = computed(() => {
  return Math.round((currentCalories / targetCalories) * 100);
});

interface MealBreakdownItem {
  name: string;
  calories: number;
  percentage: number;
  color: string;
  image: string;
}

const mealBreakdown: MealBreakdownItem[] = [
  {
    name: 'Stir Fry',
    calories: 450,
    percentage: 75,
    color: 'orange',
    image: 'https://picsum.photos/seed/stirfry/100/100',
  },
  {
    name: 'Salad',
    calories: 280,
    percentage: 60,
    color: 'green',
    image: 'https://picsum.photos/seed/salad/100/100',
  },
  {
    name: 'Smoothie',
    calories: 320,
    percentage: 85,
    color: 'purple',
    image: 'https://picsum.photos/seed/smoothie/100/100',
  },
];
</script>

<style scoped lang="scss">
.calorie-tracker {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.main-progress {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.progress-text {
  position: absolute;
  text-align: center;
}

.meal-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  flex-shrink: 0;
}

.meal-item-image {
  position: absolute;
  top: 8px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
