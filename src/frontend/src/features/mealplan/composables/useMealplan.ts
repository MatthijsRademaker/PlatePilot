import { storeToRefs } from 'pinia';
import { useMealplanStore } from '@features/mealplan/store/mealplanStore';

export function useMealplan() {
  const store = useMealplanStore();
  const {
    currentWeek,
    loading,
    error,
    totalMealsPlanned,
    suggestedRecipeIds,
    suggestionsLoading,
    suggestionsError,
  } = storeToRefs(store);

  return {
    currentWeek,
    loading,
    error,
    totalMealsPlanned,
    suggestedRecipeIds,
    suggestionsLoading,
    suggestionsError,
    setRecipeForSlot: store.setRecipeForSlot,
    clearSlot: store.clearSlot,
    navigateWeek: store.navigateWeek,
    goToCurrentWeek: store.goToCurrentWeek,
    clearWeek: store.clearWeek,
    fetchSuggestions: store.fetchSuggestions,
    clearSuggestions: store.clearSuggestions,
  };
}
