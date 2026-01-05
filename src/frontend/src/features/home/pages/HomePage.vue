<template>
  <q-page class="home-page">
    <!-- Simple Header -->
    <div class="page-header tw-flex tw-items-center tw-gap-3 tw-px-4 tw-py-3">
      <div class="header-logo">
        <q-icon name="restaurant_menu" size="20px" color="white" />
      </div>
      <span class="tw-text-lg tw-font-semibold tw-text-gray-800">PlatePilot</span>
    </div>

    <div class="tw-px-4 tw-pb-24">
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
