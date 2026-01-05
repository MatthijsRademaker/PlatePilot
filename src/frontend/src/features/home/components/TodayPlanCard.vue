<template>
  <div class="today-plan-card">
    <!-- Card Header -->
    <div class="card-header tw-px-4 tw-pt-4 tw-pb-3">
      <div class="tw-flex tw-items-center tw-gap-2">
        <q-icon name="restaurant" size="18px" color="white" class="tw-opacity-80" />
        <span class="tw-text-white tw-font-medium">Your Meal Plan Today</span>
      </div>
    </div>

    <!-- Featured Meal Content -->
    <div class="card-content tw-px-4 tw-pb-4">
      <div class="featured-meal-wrapper">
        <template v-if="featuredMeal?.recipe">
          <!-- Recipe Image -->
          <div class="meal-image-container">
            <img
              :src="getRecipeImage(featuredMeal.recipe.name)"
              :alt="featuredMeal.recipe.name"
              class="meal-image"
            />
          </div>

          <!-- Recipe Info -->
          <div class="meal-info tw-p-3">
            <div class="tw-text-sm tw-text-gray-500 tw-mb-1">
              {{ formatMealType(featuredMeal.mealType) }}
            </div>
            <div class="tw-text-lg tw-font-semibold tw-text-gray-800 tw-mb-3">
              {{ featuredMeal.recipe.name }}
            </div>

            <!-- Action Buttons -->
            <div class="tw-flex tw-flex-col tw-gap-2">
              <q-btn
                label="View Recipe"
                unelevated
                class="view-recipe-btn"
                @click="viewRecipe"
              />
              <q-btn
                label="Meal Planner Overview"
                flat
                class="overview-btn"
                @click="goToMealPlanner"
              />
            </div>
          </div>
        </template>

        <template v-else>
          <!-- Empty State -->
          <div class="empty-state tw-p-6 tw-text-center">
            <div class="empty-icon tw-mx-auto tw-mb-3">
              <q-icon name="event_busy" size="32px" color="grey-5" />
            </div>
            <div class="tw-text-gray-600 tw-mb-3">No meals planned for today</div>
            <q-btn
              label="Plan Your Meals"
              unelevated
              class="view-recipe-btn"
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
  // Generate a placeholder image based on recipe name
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
.today-plan-card {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 20px;
  overflow: hidden;
}

.featured-meal-wrapper {
  background: white;
  border-radius: 16px;
  overflow: hidden;
}

.meal-image-container {
  width: 100%;
  height: 160px;
  overflow: hidden;
}

.meal-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.view-recipe-btn {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
  border-radius: 10px;
  font-weight: 500;
  text-transform: none;
}

.overview-btn {
  color: #ff7f50 !important;
  font-weight: 500;
  text-transform: none;
}

.empty-state {
  background: white;
  border-radius: 16px;
}

.empty-icon {
  width: 64px;
  height: 64px;
  background: #f5f5f5;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
