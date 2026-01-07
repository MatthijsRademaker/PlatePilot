<template>
  <q-page class="search-page">
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-gap-3 tw-mb-4">
        <div class="header-icon">
          <q-icon name="search" size="22px" color="white" />
        </div>
        <h1 class="page-title">Search Recipes</h1>
      </div>

      <q-input
        v-model="searchQuery"
        outlined
        placeholder="What are you looking for?"
        debounce="300"
        clearable
        class="search-input"
        bg-color="white"
      >
        <template #prepend>
          <q-icon name="search" color="grey-6" />
        </template>
      </q-input>
    </header>

    <div class="tw-px-4 tw-pb-24">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <q-spinner size="40px" color="deep-orange" />
      </div>

      <!-- No Results -->
      <div v-else-if="results.length === 0 && searchQuery" class="empty-state">
        <div class="empty-icon">
          <q-icon name="search_off" size="40px" color="grey-4" />
        </div>
        <h3 class="empty-title">No recipes found</h3>
        <p class="empty-text">Try a different search term</p>
      </div>

      <!-- Initial State -->
      <div v-else-if="results.length === 0" class="empty-state">
        <div class="empty-icon empty-icon--warm">
          <q-icon name="restaurant_menu" size="40px" color="deep-orange-3" />
        </div>
        <h3 class="empty-title">Start searching</h3>
        <p class="empty-text">Find recipes by name or description</p>
      </div>

      <!-- Results Grid -->
      <div v-else class="results-grid">
        <template v-for="recipe in results" :key="recipe.id">
          <RecipeCard :recipe="recipe" @click="goToRecipe(recipe)"
        /></template>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import RecipeCard from '@features/recipe/components/RecipeCard.vue';
import { useRecipeStore } from '@features/recipe/store/recipeStore';
import type { HandlerRecipeJSON } from '@/api/generated/models';

const router = useRouter();
const recipeStore = useRecipeStore();

const searchQuery = ref('');
const results = ref<HandlerRecipeJSON[]>([]);
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
        (r.name?.toLowerCase().includes(query.toLowerCase()) ?? false) ||
        (r.description?.toLowerCase().includes(query.toLowerCase()) ?? false),
    );
  } finally {
    loading.value = false;
  }
});

function goToRecipe(recipe: HandlerRecipeJSON) {
  void router.push({ name: 'recipe-detail', params: { id: recipe.id } });
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600&display=swap');

.search-page {
  padding-top: env(safe-area-inset-top);
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 20px 16px 24px;
  color: white;
}

.header-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-title {
  font-family: 'Fraunces', serif;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
  letter-spacing: -0.3px;
}

.search-input {
  :deep(.q-field__control) {
    border-radius: 14px;
    height: 52px;
  }

  :deep(.q-field__native) {
    font-family: 'DM Sans', sans-serif;
    font-size: 16px;
    color: #2d1f1a;

    &::placeholder {
      color: #a8a0a0;
    }
  }
}

// States
.loading-state {
  display: flex;
  justify-content: center;
  padding: 60px 0;
}

.empty-state {
  background: white;
  border-radius: 24px;
  padding: 48px 24px;
  text-align: center;
  margin-top: 16px;
  box-shadow: 0 4px 20px rgba(45, 31, 26, 0.04);
  border: 1px solid rgba(45, 31, 26, 0.04);
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: #f5f2f0;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;

  &--warm {
    background: #fff5f2;
  }
}

.empty-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0 0 8px;
}

.empty-text {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  color: #a8a0a0;
  margin: 0;
}

// Results Grid
.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-top: 16px;

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}
</style>
