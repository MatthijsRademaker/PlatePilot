<template>
  <q-page padding>
    <h1 class="text-h4 q-mt-none q-mb-md">Search Recipes</h1>

    <q-input
      v-model="searchQuery"
      outlined
      placeholder="Search for recipes..."
      debounce="300"
      clearable
    >
      <template #prepend>
        <q-icon name="search" />
      </template>
    </q-input>

    <div class="q-mt-lg">
      <div v-if="loading" class="row justify-center">
        <q-spinner size="lg" color="primary" />
      </div>

      <div v-else-if="results.length === 0 && searchQuery" class="text-center text-grey">
        No recipes found for "{{ searchQuery }}"
      </div>

      <div v-else-if="results.length === 0" class="text-center text-grey">
        Start typing to search for recipes
      </div>

      <div v-else class="row q-col-gutter-md">
        <div
          v-for="recipe in results"
          :key="recipe.id"
          class="col-12 col-sm-6 col-md-4 col-lg-3"
        >
          <RecipeCard :recipe="recipe" @click="goToRecipe(recipe)" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import RecipeCard from '@features/recipe/components/RecipeCard.vue';
import { useRecipeStore } from '@features/recipe/store/recipeStore';
import type { Recipe } from '@features/recipe/types/recipe';

const router = useRouter();
const recipeStore = useRecipeStore();

const searchQuery = ref('');
const results = ref<Recipe[]>([]);
const loading = ref(false);

watch(searchQuery, async (query) => {
  if (!query) {
    results.value = [];
    return;
  }

  loading.value = true;
  try {
    // For now, filter from loaded recipes
    // TODO: Implement server-side search
    if (recipeStore.recipes.length === 0) {
      await recipeStore.fetchRecipes();
    }
    results.value = recipeStore.recipes.filter(
      (r) =>
        r.name.toLowerCase().includes(query.toLowerCase()) ||
        r.description.toLowerCase().includes(query.toLowerCase())
    );
  } finally {
    loading.value = false;
  }
});

function goToRecipe(recipe: Recipe) {
  void router.push({ name: 'recipe-detail', params: { id: recipe.id } });
}
</script>
