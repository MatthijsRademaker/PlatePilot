<template>
  <q-page class="recipe-detail-page">
    <div v-if="loading" class="row justify-center q-pa-xl">
      <q-spinner size="lg" color="primary" />
    </div>

    <div v-else-if="error" class="text-center q-pa-xl">
      <div class="error-icon tw-mx-auto tw-mb-4">
        <q-icon name="error" size="48px" color="negative" />
      </div>
      <p class="text-negative tw-mb-4">{{ error }}</p>
      <q-btn color="primary" label="Go Back" unelevated rounded @click="router.back()" />
    </div>

    <template v-else-if="recipe">
      <!-- Hero Section -->
      <div class="recipe-hero">
        <q-btn
          flat
          round
          icon="arrow_back"
          class="back-btn"
          @click="router.back()"
        />
        <div class="hero-content">
          <q-icon name="restaurant_menu" size="64px" color="white" />
        </div>
      </div>

      <div class="content-wrapper tw-px-4 tw-pb-8">
        <!-- Recipe Info Card -->
        <q-card flat class="recipe-info-card tw-mb-4">
          <q-card-section>
            <h1 class="text-h5 q-mt-none tw-font-bold">{{ recipe.name }}</h1>
            <p class="text-body1 text-grey-7">{{ recipe.description }}</p>

            <div class="row q-gutter-sm q-mt-md">
              <q-chip v-if="recipe.prepTime" icon="schedule" class="info-chip">
                Prep: {{ recipe.prepTime }}
              </q-chip>
              <q-chip v-if="recipe.cookTime" icon="local_fire_department" class="info-chip info-chip--cook">
                Cook: {{ recipe.cookTime }}
              </q-chip>
            </div>

            <div class="row q-gutter-sm q-mt-md">
              <q-chip v-if="recipe.cuisine" class="tag-chip">
                {{ recipe.cuisine.name }}
              </q-chip>
              <q-chip v-if="recipe.mainIngredient" class="tag-chip tag-chip--secondary">
                {{ recipe.mainIngredient.name }}
              </q-chip>
            </div>
          </q-card-section>
        </q-card>

        <!-- Ingredients Card -->
        <q-card flat class="section-card tw-mb-4">
          <q-card-section>
            <div class="tw-flex tw-items-center tw-gap-2 tw-mb-3">
              <div class="section-icon">
                <q-icon name="checklist" size="18px" color="white" />
              </div>
              <div class="text-h6 tw-font-semibold">Ingredients</div>
            </div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="ingredient in recipe.ingredients" :key="ingredient.id ?? ''">
              <q-item-section avatar>
                <q-icon name="check_circle" color="primary" />
              </q-item-section>
              <q-item-section>{{ ingredient.name }}</q-item-section>
            </q-item>
          </q-list>
        </q-card>

        <!-- Directions Card -->
        <q-card flat class="section-card">
          <q-card-section>
            <div class="tw-flex tw-items-center tw-gap-2 tw-mb-3">
              <div class="section-icon section-icon--accent">
                <q-icon name="format_list_numbered" size="18px" color="white" />
              </div>
              <div class="text-h6 tw-font-semibold">Directions</div>
            </div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="(direction, index) in recipe.directions" :key="index">
              <q-item-section avatar>
                <q-avatar class="step-avatar" size="sm">
                  {{ index + 1 }}
                </q-avatar>
              </q-item-section>
              <q-item-section>{{ direction }}</q-item-section>
            </q-item>
          </q-list>
        </q-card>
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
</script>

<style scoped lang="scss">
.recipe-detail-page {
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.recipe-hero {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  height: 200px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-btn {
  position: absolute;
  top: 12px;
  left: 12px;
  color: white;
  background: rgba(255, 255, 255, 0.2);
}

.hero-content {
  width: 100px;
  height: 100px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content-wrapper {
  margin-top: -40px;
  position: relative;
}

.recipe-info-card,
.section-card {
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
}

.info-chip {
  background: #fff5f2 !important;
  color: #ff7f50 !important;

  &--cook {
    background: #fff0eb !important;
    color: #ff6347 !important;
  }
}

.tag-chip {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;

  &--secondary {
    background: linear-gradient(135deg, #ffa07a 0%, #ff7f50 100%) !important;
  }
}

.section-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;

  &--accent {
    background: linear-gradient(135deg, #ff6347 0%, #e55039 100%);
  }
}

.step-avatar {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
  font-weight: 600;
}

.error-icon {
  width: 80px;
  height: 80px;
  background: #ffebee;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
