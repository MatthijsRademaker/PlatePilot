<template>
  <div class="recipe-list">
    <div v-if="loading && recipes.length === 0" class="row justify-center q-pa-lg">
      <q-spinner size="lg" color="primary" />
    </div>

    <div v-else-if="error" class="text-center q-pa-lg">
      <q-icon name="error" size="xl" color="negative" />
      <p class="text-negative">{{ error }}</p>
      <q-btn color="primary" label="Retry" @click="$emit('refresh')" />
    </div>

    <div v-else-if="recipes.length === 0" class="text-center q-pa-lg">
      <q-icon name="restaurant_menu" size="xl" color="grey" />
      <p class="text-grey">No recipes found</p>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        class="col-12 col-sm-6 col-md-4 col-lg-3"
      >
        <RecipeCard :recipe="recipe" @click="$emit('select', recipe)" />
      </div>
    </div>

    <div v-if="totalPages > 1" class="row justify-center q-mt-lg">
      <q-pagination
        :model-value="pageIndex"
        :max="totalPages"
        direction-links
        boundary-links
        @update:model-value="$emit('page', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Recipe } from '../types';
import RecipeCard from './RecipeCard.vue';

defineProps<{
  recipes: Recipe[];
  loading: boolean;
  error: string | null;
  pageIndex: number;
  totalPages: number;
}>();

defineEmits<{
  select: [recipe: Recipe];
  page: [page: number];
  refresh: [];
}>();
</script>
