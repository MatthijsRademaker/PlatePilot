<template>
  <q-page class="home-page">
    <div class="tw-px-4 tw-py-4 tw-pb-24">
      <!-- Today's Meal Plan Card -->
      <TodayPlanCard class="tw-mb-6" />

      <!-- Daily Calorie Tracker -->
      <DailyCalorieTracker class="tw-mb-6" />

      <!-- Recipe Suggestions -->
      <RecipeSuggestions />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import TodayPlanCard from '@features/home/components/TodayPlanCard.vue';
import DailyCalorieTracker from '@features/home/components/DailyCalorieTracker.vue';
import RecipeSuggestions from '@features/home/components/RecipeSuggestions.vue';
</script>

<style scoped lang="scss">
.home-page {
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: transparent;
}

.header-logo {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
