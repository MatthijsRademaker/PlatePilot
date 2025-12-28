import { apiClient } from '@/shared/api';
import type { PaginatedResponse } from '@/shared/types';
import type { Recipe, CreateRecipeRequest } from '../types';

const BASE_PATH = '/v1/recipe';

export const recipeApi = {
  async getById(id: string): Promise<Recipe> {
    return apiClient.get<Recipe>(`${BASE_PATH}/${id}`);
  },

  async getAll(pageIndex = 1, pageSize = 20): Promise<PaginatedResponse<Recipe>> {
    return apiClient.get<PaginatedResponse<Recipe>>(`${BASE_PATH}/all`, {
      params: { pageIndex, pageSize },
    });
  },

  async getSimilar(recipeId: string, amount = 5): Promise<Recipe[]> {
    return apiClient.get<Recipe[]>(`${BASE_PATH}/similar`, {
      params: { recipe: recipeId, amount },
    });
  },

  async getByCuisine(cuisineId: string): Promise<Recipe[]> {
    return apiClient.get<Recipe[]>(`${BASE_PATH}/cuisine/${cuisineId}`);
  },

  async getByIngredient(ingredientId: string): Promise<Recipe[]> {
    return apiClient.get<Recipe[]>(`${BASE_PATH}/ingredient/${ingredientId}`);
  },

  async getByAllergy(allergyId: string): Promise<Recipe[]> {
    return apiClient.get<Recipe[]>(`${BASE_PATH}/allergy/${allergyId}`);
  },

  async create(request: CreateRecipeRequest): Promise<Recipe> {
    return apiClient.post<Recipe>(`${BASE_PATH}/create`, request);
  },
};
