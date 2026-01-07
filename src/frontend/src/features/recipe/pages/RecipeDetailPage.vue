<template>
  <q-page class="recipe-detail-page">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <q-spinner size="48px" color="deep-orange" />
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">
        <q-icon name="error_outline" size="48px" color="red-5" />
      </div>
      <p class="error-text">{{ error }}</p>
      <q-btn
        label="Go Back"
        unelevated
        no-caps
        class="error-btn"
        @click="router.back()"
      />
    </div>

    <!-- Recipe Content -->
    <template v-else-if="recipe">
      <!-- Hero Section -->
      <div class="recipe-hero">
        <img :src="getRecipeImage(recipe.name)" :alt="recipe.name" class="hero-image" />
        <div class="hero-overlay"></div>
        <q-btn
          flat
          round
          icon="arrow_back"
          class="back-btn"
          @click="router.back()"
        />
        <div class="hero-content">
          <div v-if="recipe.cuisine" class="cuisine-badge">
            {{ recipe.cuisine.name }}
          </div>
          <h1 class="recipe-title">{{ recipe.name }}</h1>
          <div v-if="recipe.prepTime || recipe.cookTime" class="time-badges">
            <span v-if="recipe.prepTime" class="time-badge">
              <q-icon name="schedule" size="14px" />
              Prep: {{ recipe.prepTime }}
            </span>
            <span v-if="recipe.cookTime" class="time-badge">
              <q-icon name="local_fire_department" size="14px" />
              Cook: {{ recipe.cookTime }}
            </span>
          </div>
        </div>
      </div>

      <div class="content-wrapper">
        <!-- Description -->
        <p v-if="recipe.description" class="recipe-description">
          {{ recipe.description }}
        </p>

        <!-- Ingredients Card -->
        <div class="section-card">
          <div class="section-header">
            <div class="section-icon">
              <q-icon name="checklist" size="18px" color="white" />
            </div>
            <h2 class="section-title">Ingredients</h2>
          </div>
          <ul class="ingredients-list">
            <li v-for="ingredient in recipe.ingredients" :key="ingredient.id ?? ''">
              <q-icon name="check_circle" size="18px" class="check-icon" />
              <span>{{ ingredient.name }}</span>
            </li>
          </ul>
        </div>

        <!-- Directions Card -->
        <div class="section-card">
          <div class="section-header">
            <div class="section-icon section-icon--accent">
              <q-icon name="format_list_numbered" size="18px" color="white" />
            </div>
            <h2 class="section-title">Directions</h2>
          </div>
          <ol class="directions-list">
            <li v-for="(direction, index) in recipe.directions" :key="index">
              <span class="step-number">{{ index + 1 }}</span>
              <span class="step-text">{{ direction }}</span>
            </li>
          </ol>
        </div>
      </div>
    </template>
  </q-page>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
import { computed } from 'vue';
import { useRecipeDetail } from '@features/recipe/composables/useRecipe';

const router = useRouter();
const route = useRoute();

const recipeId = computed(() => route.params.id as string);
const { recipe, loading, error } = useRecipeDetail(() => recipeId.value);

function getRecipeImage(recipeName: string | undefined): string {
  const seed = recipeName?.replace(/\s+/g, '-').toLowerCase() || 'default';
  return `https://picsum.photos/seed/${seed}/800/500`;
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,500;9..144,600;9..144,700&display=swap');

.recipe-detail-page {
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

// Loading & Error States
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 24px;
  text-align: center;
}

.error-icon {
  width: 80px;
  height: 80px;
  background: #ffebee;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.error-text {
  font-family: 'DM Sans', sans-serif;
  font-size: 16px;
  color: #c62828;
  margin: 0 0 20px;
}

.error-btn {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
  border-radius: 12px;
  font-family: 'DM Sans', sans-serif;
  font-weight: 600;
  padding: 0 24px;
  height: 44px;
}

// Hero Section
.recipe-hero {
  position: relative;
  height: 280px;
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.1) 0%, rgba(0, 0, 0, 0.6) 100%);
}

.back-btn {
  position: absolute;
  top: calc(env(safe-area-inset-top) + 12px);
  left: 12px;
  color: white;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
}

.hero-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20px 16px;
  color: white;
}

.cuisine-badge {
  display: inline-block;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 6px 14px;
  border-radius: 20px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 12px;
}

.recipe-title {
  font-family: 'Fraunces', serif;
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 12px;
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.time-badges {
  display: flex;
  gap: 10px;
}

.time-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(8px);
  padding: 6px 12px;
  border-radius: 8px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 500;
}

// Content
.content-wrapper {
  padding: 20px 16px 100px;
}

.recipe-description {
  font-family: 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 400;
  color: #6b5f5a;
  line-height: 1.6;
  margin: 0 0 24px;
}

// Section Cards
.section-card {
  background: white;
  border-radius: 20px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 4px 20px rgba(45, 31, 26, 0.06);
  border: 1px solid rgba(45, 31, 26, 0.04);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.section-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;

  &--accent {
    background: linear-gradient(135deg, #ff6347 0%, #e5533d 100%);
  }
}

.section-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
}

// Ingredients List
.ingredients-list {
  list-style: none;
  padding: 0;
  margin: 0;

  li {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid #f5f2f0;
    font-family: 'DM Sans', sans-serif;
    font-size: 15px;
    color: #4a3f3a;

    &:last-child {
      border-bottom: none;
    }
  }

  .check-icon {
    color: #ff6347;
    flex-shrink: 0;
  }
}

// Directions List
.directions-list {
  list-style: none;
  padding: 0;
  margin: 0;
  counter-reset: step;

  li {
    display: flex;
    gap: 14px;
    padding: 14px 0;
    border-bottom: 1px solid #f5f2f0;

    &:last-child {
      border-bottom: none;
    }
  }

  .step-number {
    width: 28px;
    height: 28px;
    background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: 'DM Sans', sans-serif;
    font-size: 13px;
    font-weight: 600;
    color: white;
    flex-shrink: 0;
  }

  .step-text {
    font-family: 'DM Sans', sans-serif;
    font-size: 15px;
    color: #4a3f3a;
    line-height: 1.5;
    padding-top: 4px;
  }
}
</style>
