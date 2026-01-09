import customInstance from '@/api/mutator/custom-instance';
import type { WeekPlan } from '@features/mealplan/types/mealplan';

export interface WeekPlanSavePayload {
  startDate: string;
  endDate: string;
  days: DayPlanSavePayload[];
}

export interface DayPlanSavePayload {
  date: string;
  meals: MealSlotSavePayload[];
}

export interface MealSlotSavePayload {
  mealType: string;
  recipeId?: string;
}

export function getWeekPlan(startDate: string): Promise<WeekPlan> {
  return customInstance({
    url: '/mealplan/week',
    method: 'GET',
    params: { startDate },
  });
}

export function upsertWeekPlan(payload: WeekPlanSavePayload): Promise<WeekPlan> {
  return customInstance({
    url: '/mealplan/week',
    method: 'PUT',
    data: payload,
  });
}
