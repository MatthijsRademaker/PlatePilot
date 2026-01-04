<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
      <h1 class="text-h4 q-ma-none">Recipes</h1>
      <q-btn color="primary" icon="refresh" flat round @click="refresh" />
    </div>

    <RecipeList
      :recipes="recipes"
      :loading="loading"
      :error="error"
      :page-index="pageIndex"
      :total-pages="totalPages"
      @select="goToRecipe"
      @page="loadPage"
      @refresh="refresh"
    />
  </q-page>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useRecipeList } from '@features/recipe/composables/useRecipe';
import RecipeList from '@features/recipe/components/RecipeList.vue';
import type { Recipe } from '@features/recipe/types/recipe';

const router = useRouter();
const { recipes, loading, error, pageIndex, totalPages, loadPage, refresh } = useRecipeList();

function goToRecipe(recipe: Recipe) {
  void router.push({ name: 'recipe-detail', params: { id: recipe.id } });
}
</script>
