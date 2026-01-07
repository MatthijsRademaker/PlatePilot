<template>
  <div class="today-plan-card">
    <!-- Card Header -->
    <div class="card-header">
      <div class="tw-flex tw-items-center tw-gap-2">
        <div class="header-icon">
          <q-icon name="restaurant" size="16px" color="white" />
        </div>
        <span class="header-title">Your Meal Plan Today</span>
      </div>
    </div>

    <!-- Featured Meal Content -->
    <div class="card-content">
      <div class="featured-meal-wrapper">
        <template v-if="featuredMeal?.recipe">
          <!-- Recipe Image -->
          <div class="meal-image-container">
            <img
              :src="getRecipeImage(featuredMeal.recipe.name)"
              :alt="featuredMeal.recipe.name"
              class="meal-image"
            />
            <div class="meal-badge">
              {{ formatMealType(featuredMeal.mealType) }}
            </div>
          </div>

          <!-- Recipe Info -->
          <div class="meal-info">
            <h3 class="meal-name">
              {{ featuredMeal.recipe.name }}
            </h3>

            <!-- Action Buttons -->
            <div class="tw-flex tw-gap-3 tw-mt-4">
              <q-btn
                label="View Recipe"
                unelevated
                no-caps
                class="view-btn tw-flex-1"
                @click="viewRecipe"
              />
              <q-btn
                icon="mdi-calendar"
                flat
                round
                class="plan-btn"
                @click="goToMealPlanner"
              />
            </div>
          </div>
        </template>

        <template v-else>
          <!-- Empty State -->
          <div class="empty-state">
            <div class="empty-icon">
              <q-icon name="event_busy" size="28px" color="grey-5" />
            </div>
            <p class="empty-text">No meals planned for today</p>
            <q-btn
              label="Plan Your Meals"
              unelevated
              no-caps
              class="view-btn"
              @click="goToMealPlanner"
            />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useMealplanStore } from '@features/mealplan/store/mealplanStore';
import type { MealSlot, MealType } from '@features/mealplan/types/mealplan';

const router = useRouter();
const mealplanStore = useMealplanStore();

const todayDateStr = computed(() => {
  return new Date().toISOString().split('T')[0];
});

const todayPlan = computed(() => {
  return mealplanStore.currentWeek.days.find((day) => day.date === todayDateStr.value);
});

const featuredMeal = computed((): MealSlot | null => {
  if (!todayPlan.value) return null;

  // Find first meal with a recipe (prioritize dinner > lunch > breakfast)
  const priority: MealType[] = ['dinner', 'lunch', 'breakfast'];
  for (const mealType of priority) {
    const meal = todayPlan.value.meals.find(
      (m) => m.mealType === mealType && m.recipe !== null,
    );
    if (meal) return meal;
  }
  return null;
});

function formatMealType(mealType: MealType): string {
  return mealType.charAt(0).toUpperCase() + mealType.slice(1);
}

function getRecipeImage(recipeName: string | undefined): string {
  const seed = recipeName?.replace(/\s+/g, '-').toLowerCase() || 'default';
  return `https://picsum.photos/seed/${seed}/400/240`;
}

function viewRecipe() {
  if (featuredMeal.value?.recipe?.id) {
    void router.push({ name: 'recipe-detail', params: { id: featuredMeal.value.recipe.id } });
  }
}

function goToMealPlanner() {
  void router.push({ name: 'mealplan' });
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700&family=Fraunces:opsz,wght@9..144,600&display=swap');

.today-plan-card {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(255, 99, 71, 0.25);
}

.card-header {
  padding: 16px 16px 12px;
}

.header-icon {
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.card-content {
  padding: 0 12px 12px;
}

.featured-meal-wrapper {
  background: white;
  border-radius: 18px;
  overflow: hidden;
}

.meal-image-container {
  position: relative;
  width: 100%;
  height: 140px;
  overflow: hidden;
}

.meal-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.meal-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: rgba(45, 31, 26, 0.75);
  backdrop-filter: blur(8px);
  padding: 6px 12px;
  border-radius: 20px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: white;
}

.meal-info {
  padding: 16px;
}

.meal-name {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
  line-height: 1.3;
}

.view-btn {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
  border-radius: 12px;
  font-family: 'DM Sans', sans-serif;
  font-weight: 600;
  height: 44px;
}

.plan-btn {
  background: #fff5f2 !important;
  color: #ff6347 !important;
  width: 44px;
  height: 44px;
}

.empty-state {
  padding: 32px 24px;
  text-align: center;
}

.empty-icon {
  width: 56px;
  height: 56px;
  background: #f5f2f0;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
}

.empty-text {
  font-family: 'DM Sans', sans-serif;
  font-size: 15px;
  color: #6b5f5a;
  margin: 0 0 16px;
}
</style>
