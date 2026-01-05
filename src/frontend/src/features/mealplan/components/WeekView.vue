<template>
  <div class="week-view">
    <div class="week-nav tw-flex tw-items-center tw-justify-between tw-mb-4">
      <q-btn flat round icon="chevron_left" class="nav-btn" @click="$emit('prev')" />
      <div class="text-h6 tw-font-semibold">
        {{ formatDateRange(weekPlan.startDate, weekPlan.endDate) }}
      </div>
      <q-btn flat round icon="chevron_right" class="nav-btn" @click="$emit('next')" />
    </div>

    <div class="row q-col-gutter-sm">
      <div
        v-for="day in weekPlan.days"
        :key="day.date"
        class="col"
      >
        <q-card flat class="day-card">
          <q-card-section class="day-header q-pa-sm text-center">
            <div class="text-caption text-weight-medium tw-text-white">
              {{ formatDayName(day.date) }}
            </div>
            <div class="text-caption tw-text-white/80">
              {{ formatDate(day.date) }}
            </div>
          </q-card-section>

          <q-card-section class="q-pa-sm q-gutter-sm">
            <MealSlotCard
              v-for="meal in day.meals"
              :key="meal.id"
              :meal-slot="meal"
              @click="$emit('slot-click', meal)"
              @clear="$emit('slot-clear', meal)"
            />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { WeekPlan, MealSlot } from '@features/mealplan/types/mealplan';
import MealSlotCard from '@features/mealplan/components/MealSlotCard.vue';

defineProps<{
  weekPlan: WeekPlan;
}>();

defineEmits<{
  prev: [];
  next: [];
  'slot-click': [slot: MealSlot];
  'slot-clear': [slot: MealSlot];
}>();

function formatDayName(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', { weekday: 'short' });
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

function formatDateRange(startStr: string, endStr: string): string {
  const start = new Date(startStr);
  const end = new Date(endStr);
  const startMonth = start.toLocaleDateString('en-US', { month: 'short' });
  const endMonth = end.toLocaleDateString('en-US', { month: 'short' });

  if (startMonth === endMonth) {
    return `${startMonth} ${start.getDate()} - ${end.getDate()}, ${start.getFullYear()}`;
  }
  return `${startMonth} ${start.getDate()} - ${endMonth} ${end.getDate()}, ${start.getFullYear()}`;
}
</script>

<style scoped lang="scss">
.week-view {
  overflow-x: auto;
}

.nav-btn {
  color: #ff7f50;
  background: #fff5f2;

  &:hover {
    background: #ffebe5;
  }
}

.day-card {
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.day-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
}
</style>
