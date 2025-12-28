<template>
  <div class="week-view">
    <div class="row items-center justify-between q-mb-md">
      <q-btn flat icon="chevron_left" @click="$emit('prev')" />
      <div class="text-h6">
        {{ formatDateRange(weekPlan.startDate, weekPlan.endDate) }}
      </div>
      <q-btn flat icon="chevron_right" @click="$emit('next')" />
    </div>

    <div class="row q-col-gutter-sm">
      <div
        v-for="day in weekPlan.days"
        :key="day.date"
        class="col"
      >
        <q-card flat bordered>
          <q-card-section class="q-pa-sm text-center bg-grey-2">
            <div class="text-caption text-weight-medium">
              {{ formatDayName(day.date) }}
            </div>
            <div class="text-caption text-grey">
              {{ formatDate(day.date) }}
            </div>
          </q-card-section>

          <q-card-section class="q-pa-sm q-gutter-sm">
            <MealSlotCard
              v-for="meal in day.meals"
              :key="meal.id"
              :slot="meal"
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
import type { WeekPlan, MealSlot } from '../types';
import MealSlotCard from './MealSlotCard.vue';

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
</style>
