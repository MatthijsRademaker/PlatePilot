<template>
  <q-card class="today-plan-card">
    <q-card-section>
      <div class="row items-center q-mb-md">
        <q-icon name="today" size="24px" color="primary" class="q-mr-sm" />
        <div class="text-h6">My Plan For Today</div>
      </div>

      <template v-if="hasMealsPlanned">
        <div class="column q-gutter-sm">
          <div
            v-for="meal in todayMeals"
            :key="meal.mealType"
            class="meal-row row items-center q-pa-sm rounded-borders cursor-pointer"
            :class="{ 'meal-row--clickable': meal.recipe }"
            @click="handleMealClick(meal)"
          >
            <q-icon
              :name="getMealIcon(meal.mealType)"
              size="20px"
              :color="meal.recipe ? 'primary' : 'grey'"
              class="q-mr-sm"
            />
            <div class="column">
              <span class="text-caption text-weight-medium text-uppercase">
                {{ formatMealType(meal.mealType) }}
              </span>
              <span
                class="text-body2"
                :class="meal.recipe ? '' : 'text-grey'"
              >
                {{ meal.recipe?.name || 'Not planned' }}
              </span>
            </div>
            <q-icon
              v-if="meal.recipe"
              name="chevron_right"
              size="20px"
              color="grey"
              class="q-ml-auto"
            />
          </div>
        </div>
      </template>

      <template v-else>
        <div class="text-center q-py-md">
          <q-icon name="event_busy" size="48px" color="grey-5" />
          <div class="text-body1 text-grey q-mt-sm">No meals planned for today</div>
          <q-btn
            label="Plan Your Meals"
            color="primary"
            outline
            class="q-mt-md"
            @click="navigateToMealPlanner"
          />
        </div>
      </template>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useMealplanStore, type MealSlot, type MealType } from '@/features/mealplan';

const router = useRouter();
const mealplanStore = useMealplanStore();

const todayDateStr = computed(() => {
  return new Date().toISOString().split('T')[0];
});

const todayPlan = computed(() => {
  return mealplanStore.currentWeek.days.find(
    (day) => day.date === todayDateStr.value
  );
});

const todayMeals = computed((): MealSlot[] => {
  if (!todayPlan.value) {
    return [];
  }
  // Only show breakfast, lunch, dinner (not snack)
  return todayPlan.value.meals.filter(
    (meal) => ['breakfast', 'lunch', 'dinner'].includes(meal.mealType)
  );
});

const hasMealsPlanned = computed(() => {
  return todayMeals.value.some((meal) => meal.recipe !== null);
});

function getMealIcon(mealType: MealType): string {
  const icons: Record<MealType, string> = {
    breakfast: 'free_breakfast',
    lunch: 'lunch_dining',
    dinner: 'dinner_dining',
    snack: 'cookie',
  };
  return icons[mealType] || 'restaurant';
}

function formatMealType(mealType: MealType): string {
  return mealType.charAt(0).toUpperCase() + mealType.slice(1);
}

function handleMealClick(meal: MealSlot) {
  if (meal.recipe) {
    void router.push({ name: 'recipe-detail', params: { id: meal.recipe.id } });
  }
}

function navigateToMealPlanner() {
  void router.push({ name: 'mealplan' });
}
</script>

<style scoped lang="scss">
.today-plan-card {
  width: 100%;
}

.meal-row {
  background: rgba(0, 0, 0, 0.02);
  transition: all 0.2s;

  &--clickable:hover {
    background: rgba(var(--q-primary-rgb), 0.1);
    transform: translateX(4px);
  }
}
</style>
