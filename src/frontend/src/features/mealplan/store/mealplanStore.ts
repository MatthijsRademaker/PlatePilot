import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { HandlerRecipeJSON } from '@/api/generated/models';
import type { WeekPlan, DayPlan, MealSlot, MealType } from '@features/mealplan/types/mealplan';
import { postMealplanSuggest } from '@/api/generated/platepilot';

function formatDate(date: Date): string {
  return date.toISOString().split('T')[0] as string;
}

function generateWeekPlan(startDate: Date): WeekPlan {
  const days: DayPlan[] = [];
  const mealTypes: MealType[] = ['breakfast', 'lunch', 'dinner'];

  for (let i = 0; i < 7; i++) {
    const date = new Date(startDate);
    date.setDate(date.getDate() + i);
    const dateStr = formatDate(date);

    const meals: MealSlot[] = mealTypes.map((mealType) => ({
      id: `${dateStr}-${mealType}`,
      date: dateStr,
      mealType,
      recipe: null,
    }));

    days.push({ date: dateStr, meals });
  }

  const endDate = new Date(startDate);
  endDate.setDate(endDate.getDate() + 6);

  return {
    startDate: formatDate(startDate),
    endDate: formatDate(endDate),
    days,
  };
}

function getWeekStart(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay();
  const diff = d.getDate() - day + (day === 0 ? -6 : 1);
  d.setDate(diff);
  d.setHours(0, 0, 0, 0);
  return d;
}

export const useMealplanStore = defineStore('mealplan', () => {
  // State
  const currentWeek = ref<WeekPlan>(generateWeekPlan(getWeekStart(new Date())));
  const suggestedRecipeIds = ref<string[]>([]);
  const loading = ref(false);
  const suggestionsLoading = ref(false);
  const error = ref<string | null>(null);
  const suggestionsError = ref<string | null>(null);

  // Getters
  const plannedRecipeIds = computed(() => {
    const ids: string[] = [];
    currentWeek.value.days.forEach((day) => {
      day.meals.forEach((meal) => {
        if (meal.recipe?.id) {
          ids.push(meal.recipe.id);
        }
      });
    });
    return ids;
  });

  const totalMealsPlanned = computed(() => plannedRecipeIds.value.length);

  // Actions
  function setRecipeForSlot(slotId: string, recipe: HandlerRecipeJSON | null) {
    for (const day of currentWeek.value.days) {
      const slot = day.meals.find((m) => m.id === slotId);
      if (slot) {
        slot.recipe = recipe;
        break;
      }
    }
  }

  function clearSlot(slotId: string) {
    setRecipeForSlot(slotId, null);
  }

  function navigateWeek(direction: 'prev' | 'next') {
    const currentStart = new Date(currentWeek.value.startDate);
    const offset = direction === 'next' ? 7 : -7;
    currentStart.setDate(currentStart.getDate() + offset);
    currentWeek.value = generateWeekPlan(currentStart);
  }

  function goToCurrentWeek() {
    currentWeek.value = generateWeekPlan(getWeekStart(new Date()));
  }

  function clearWeek() {
    currentWeek.value.days.forEach((day) => {
      day.meals.forEach((meal) => {
        meal.recipe = null;
      });
    });
  }

  async function fetchSuggestions(amount: number = 5): Promise<void> {
    suggestionsLoading.value = true;
    suggestionsError.value = null;
    suggestedRecipeIds.value = [];

    try {
      const result = await postMealplanSuggest({
        alreadySelectedRecipeIds: plannedRecipeIds.value,
        amount,
      });
      suggestedRecipeIds.value = result.recipeIds ?? [];
    } catch (err) {
      suggestionsError.value =
        err instanceof Error ? err.message : 'Failed to fetch suggestions';
    } finally {
      suggestionsLoading.value = false;
    }
  }

  function clearSuggestions() {
    suggestedRecipeIds.value = [];
    suggestionsError.value = null;
  }

  return {
    // State
    currentWeek,
    suggestedRecipeIds,
    loading,
    suggestionsLoading,
    error,
    suggestionsError,
    // Getters
    plannedRecipeIds,
    totalMealsPlanned,
    // Actions
    setRecipeForSlot,
    clearSlot,
    navigateWeek,
    goToCurrentWeek,
    clearWeek,
    fetchSuggestions,
    clearSuggestions,
  };
});
