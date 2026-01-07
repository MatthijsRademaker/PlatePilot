<template>
  <q-page class="home-page">
    <!-- Page Header -->
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div>
          <p class="greeting">Good {{ timeOfDay }},</p>
          <h1 class="title">What's cooking?</h1>
        </div>
        <div class="header-avatar">
          <q-icon name="person" size="24px" color="white" />
        </div>
      </div>
    </header>

    <div class="tw-px-4 tw-py-4 tw-pb-24">
      <!-- Today's Meal Plan Card -->
      <TodayPlanCard class="tw-mb-5" />

      <!-- Daily Calorie Tracker -->
      <DailyCalorieTracker class="tw-mb-5" />

      <!-- Recipe Suggestions -->
      <RecipeSuggestions />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import TodayPlanCard from '@features/home/components/TodayPlanCard.vue';
import DailyCalorieTracker from '@features/home/components/DailyCalorieTracker.vue';
import RecipeSuggestions from '@features/home/components/RecipeSuggestions.vue';

const timeOfDay = computed(() => {
  const hour = new Date().getHours();
  if (hour < 12) return 'morning';
  if (hour < 17) return 'afternoon';
  return 'evening';
});
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700&family=Fraunces:opsz,wght@9..144,400;9..144,600;9..144,700&display=swap');

.home-page {
  padding-top: env(safe-area-inset-top);
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 20px 16px 24px;
  color: white;
}

.greeting {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.85);
  margin: 0 0 4px;
}

.title {
  font-family: 'Fraunces', serif;
  font-size: 28px;
  font-weight: 600;
  letter-spacing: -0.5px;
  margin: 0;
}

.header-avatar {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
