import type { Recipe } from '@features/recipe/types/recipe';

export type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snack';

export interface MealSlot {
  id: string;
  date: string;
  mealType: MealType;
  recipe: Recipe | null;
}

export interface DayPlan {
  date: string;
  meals: MealSlot[];
}

export interface WeekPlan {
  startDate: string;
  endDate: string;
  days: DayPlan[];
}

export interface SuggestRecipesRequest {
  excludeRecipeIds?: string[];
  cuisineIds?: string[];
  excludeAllergyIds?: string[];
  amount?: number;
}
