import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { recipeApi } from '@features/recipe/api/recipeApi';
import type { Recipe } from '@features/recipe/types/recipe';

export const useRecipeStore = defineStore('recipe', () => {
  // State
  const recipes = ref<Recipe[]>([]);
  const currentRecipe = ref<Recipe | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pageIndex = ref(1);
  const pageSize = ref(20);
  const totalCount = ref(0);

  // Getters
  const hasMore = computed(() => recipes.value.length < totalCount.value);
  const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

  // Actions
  async function fetchRecipes(page = 1) {
    loading.value = true;
    error.value = null;
    try {
      const response = await recipeApi.getAll(page, pageSize.value);
      recipes.value = response.items;
      pageIndex.value = response.pageIndex;
      totalCount.value = response.totalCount;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch recipes';
    } finally {
      loading.value = false;
    }
  }

  async function fetchRecipeById(id: string) {
    loading.value = true;
    error.value = null;
    try {
      currentRecipe.value = await recipeApi.getById(id);
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch recipe';
    } finally {
      loading.value = false;
    }
  }

  async function fetchSimilarRecipes(recipeId: string, amount = 5): Promise<Recipe[]> {
    try {
      return await recipeApi.getSimilar(recipeId, amount);
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch similar recipes';
      return [];
    }
  }

  function clearCurrentRecipe() {
    currentRecipe.value = null;
  }

  return {
    // State
    recipes,
    currentRecipe,
    loading,
    error,
    pageIndex,
    pageSize,
    totalCount,
    // Getters
    hasMore,
    totalPages,
    // Actions
    fetchRecipes,
    fetchRecipeById,
    fetchSimilarRecipes,
    clearCurrentRecipe,
  };
});
