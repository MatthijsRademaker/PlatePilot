import { storeToRefs } from 'pinia';
import { onMounted, watch } from 'vue';
import { useRecipeStore } from '@features/recipe/store/recipeStore';

export function useRecipeList() {
  const store = useRecipeStore();
  const { recipes, loading, error, pageIndex, totalPages, hasMore } = storeToRefs(store);

  onMounted(() => {
    if (recipes.value.length === 0) {
      store.fetchRecipes();
    }
  });

  function loadPage(page: number) {
    store.fetchRecipes(page);
  }

  function refresh() {
    store.fetchRecipes(1);
  }

  return {
    recipes,
    loading,
    error,
    pageIndex,
    totalPages,
    hasMore,
    loadPage,
    refresh,
  };
}

export function useRecipeDetail(recipeId: () => string) {
  const store = useRecipeStore();
  const { currentRecipe, loading, error } = storeToRefs(store);

  onMounted(() => {
    store.fetchRecipeById(recipeId());
  });

  watch(recipeId, (newId) => {
    store.fetchRecipeById(newId);
  });

  return {
    recipe: currentRecipe,
    loading,
    error,
  };
}
