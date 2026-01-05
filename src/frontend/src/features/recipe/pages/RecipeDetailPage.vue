<template>
  <q-page padding>
    <div v-if="loading" class="row justify-center q-pa-xl">
      <q-spinner size="lg" color="primary" />
    </div>

    <div v-else-if="error" class="text-center q-pa-xl">
      <q-icon name="error" size="xl" color="negative" />
      <p class="text-negative">{{ error }}</p>
      <q-btn color="primary" label="Go Back" @click="router.back()" />
    </div>

    <template v-else-if="recipe">
      <q-btn flat icon="arrow_back" label="Back" class="q-mb-md" @click="router.back()" />

      <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-6">
          <div class="placeholder-image rounded-borders">
            <q-icon name="restaurant_menu" size="128px" color="grey-5" />
          </div>
        </div>

        <div class="col-12 col-md-6">
          <h1 class="text-h4 q-mt-none">{{ recipe.name }}</h1>
          <p class="text-body1">{{ recipe.description }}</p>

          <div class="row q-gutter-md q-mb-md">
            <q-chip v-if="recipe.prepTime" icon="schedule" color="primary" text-color="white">
              Prep: {{ recipe.prepTime }}
            </q-chip>
            <q-chip v-if="recipe.cookTime" icon="local_fire_department" color="orange" text-color="white">
              Cook: {{ recipe.cookTime }}
            </q-chip>
          </div>

          <div v-if="recipe.cuisine" class="q-mb-md">
            <div class="text-subtitle2 q-mb-xs">Cuisine</div>
            <q-chip color="primary" text-color="white" size="sm">
              {{ recipe.cuisine.name }}
            </q-chip>
          </div>

          <div v-if="recipe.mainIngredient" class="q-mb-md">
            <div class="text-subtitle2 q-mb-xs">Main Ingredient</div>
            <q-chip color="secondary" text-color="white" size="sm">
              {{ recipe.mainIngredient.name }}
            </q-chip>
          </div>
        </div>
      </div>

      <div class="row q-col-gutter-lg q-mt-md">
        <div class="col-12 col-md-4">
          <q-card>
            <q-card-section>
              <div class="text-h6">Ingredients</div>
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
        </div>

        <div class="col-12 col-md-8">
          <q-card>
            <q-card-section>
              <div class="text-h6">Directions</div>
            </q-card-section>
            <q-list separator>
              <q-item v-for="(direction, index) in recipe.directions" :key="index">
                <q-item-section avatar>
                  <q-avatar color="primary" text-color="white" size="sm">
                    {{ index + 1 }}
                  </q-avatar>
                </q-item-section>
                <q-item-section>{{ direction }}</q-item-section>
              </q-item>
            </q-list>
          </q-card>
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
</script>

<style scoped>
.placeholder-image {
  aspect-ratio: 16 / 9;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--q-grey-3);
}
</style>
