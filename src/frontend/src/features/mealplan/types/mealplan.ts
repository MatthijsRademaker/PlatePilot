import type { HandlerRecipeJSON } from '@/api/generated/models';

export type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snack';

export interface MealSlot {
  id: string;
  date: string;
  mealType: MealType;
  recipe: HandlerRecipeJSON | null;
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
