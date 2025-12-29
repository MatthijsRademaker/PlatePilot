import { apiClient } from '@/shared/api';
import { recipeApi } from '@/features/recipe/api';
import type { Recipe } from '@/features/recipe/types';
import type { SuggestRecipesRequest, SuggestRecipesResponse } from '../types';

const BASE_PATH = '/v1/mealplan';

export const mealplanApi = {
  /**
   * Suggests recipes based on constraints, excluding already planned recipes.
   * Fetches full recipe details after receiving IDs from the backend.
   */
  async suggestRecipes(request: SuggestRecipesRequest = {}): Promise<Recipe[]> {
    // Backend returns recipe IDs, not full recipes
    const response = await apiClient.post<SuggestRecipesResponse>(
      `${BASE_PATH}/suggest`,
      {
        alreadySelectedRecipeIds: request.excludeRecipeIds ?? [],
        amount: request.amount ?? 5,
      }
    );

    // Fetch full recipe details for each suggested ID
    if (!response.recipeIds || response.recipeIds.length === 0) {
      return [];
    }

    // Fetch all recipes in parallel
    const recipePromises = response.recipeIds.map((id) => recipeApi.getById(id));
    const recipes = await Promise.all(recipePromises);

    return recipes;
  },
};
