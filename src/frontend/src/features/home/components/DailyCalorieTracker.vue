<template>
  <div class="calorie-tracker">
    <h2 class="section-title">Daily Calorie Tracker</h2>

    <div class="tracker-content">
      <!-- Main Calorie Circle -->
      <div class="main-progress">
        <q-circular-progress
          :value="caloriePercentage"
          size="90px"
          :thickness="0.12"
          color="deep-orange"
          track-color="grey-3"
          class="main-circle"
        />
        <div class="progress-text">
          <span class="calories-value">{{ currentCalories }}</span>
          <span class="calories-target">/ {{ targetCalories }}</span>
        </div>
      </div>

      <!-- Meal Breakdown -->
      <div class="meals-grid">
        <div v-for="item in mealBreakdown" :key="item.name" class="meal-item">
          <div class="meal-progress">
            <q-circular-progress
              :value="item.percentage"
              size="52px"
              :thickness="0.1"
              :color="item.color"
              track-color="grey-2"
            />
            <div class="meal-image">
              <img :src="item.image" :alt="item.name" />
            </div>
          </div>
          <span class="meal-name">{{ item.name }}</span>
          <span class="meal-cals">{{ item.calories }} cal</span>
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
    color: 'deep-orange',
    image: 'https://picsum.photos/seed/stirfry/100/100',
  },
  {
    name: 'Salad',
    calories: 280,
    percentage: 60,
    color: 'green-6',
    image: 'https://picsum.photos/seed/salad/100/100',
  },
  {
    name: 'Smoothie',
    calories: 320,
    percentage: 85,
    color: 'purple-5',
    image: 'https://picsum.photos/seed/smoothie/100/100',
  },
];
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700&family=Fraunces:opsz,wght@9..144,600&display=swap');

.calorie-tracker {
  background: white;
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 31, 26, 0.06);
  border: 1px solid rgba(45, 31, 26, 0.04);
}

.section-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0 0 16px;
}

.tracker-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.main-progress {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.progress-text {
  position: absolute;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.calories-value {
  font-family: 'DM Sans', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #2d1f1a;
  line-height: 1;
}

.calories-target {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 500;
  color: #a8a0a0;
  margin-top: 2px;
}

.meals-grid {
  display: flex;
  gap: 16px;
  flex: 1;
  overflow-x: auto;
  padding: 4px;

  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.meal-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}

.meal-progress {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.meal-image {
  position: absolute;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.meal-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #4a3f3a;
  margin-top: 8px;
  max-width: 60px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meal-cals {
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 500;
  color: #a8a0a0;
  margin-top: 2px;
}
</style>
