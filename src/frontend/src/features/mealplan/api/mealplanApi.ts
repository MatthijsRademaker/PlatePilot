import { apiClient } from '@/shared/api';
import type { Recipe } from '@features/recipe/types/recipe';
import type { SuggestRecipesRequest } from '@features/mealplan/types/mealplan';

const BASE_PATH = '/v1/mealplan';

export const mealplanApi = {
  async suggestRecipes(request: SuggestRecipesRequest = {}): Promise<Recipe[]> {
    return apiClient.post<Recipe[]>(`${BASE_PATH}/suggest`, request);
  },
};
