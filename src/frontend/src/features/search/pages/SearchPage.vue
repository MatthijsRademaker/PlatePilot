<template>
  <q-page class="search-page">
    <div class="page-header">
      <div class="tw-flex tw-items-center tw-gap-3 tw-mb-4">
        <div class="header-icon">
          <q-icon name="search" size="24px" color="white" />
        </div>
        <h1 class="text-h5 q-ma-none tw-font-semibold">Search Recipes</h1>
      </div>

      <q-input
        v-model="searchQuery"
        outlined
        placeholder="Search for recipes..."
        debounce="300"
        clearable
        class="search-input"
        bg-color="white"
      >
        <template #prepend>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <div class="tw-px-4 tw-pb-4">
      <div v-if="loading" class="row justify-center tw-py-8">
        <q-spinner size="lg" color="primary" />
      </div>

      <div v-else-if="results.length === 0 && searchQuery" class="empty-state tw-text-center tw-py-12">
        <div class="empty-icon tw-mx-auto tw-mb-4">
          <q-icon name="search_off" size="48px" color="grey-5" />
        </div>
        <div class="text-body1 text-grey-6">No recipes found for "{{ searchQuery }}"</div>
        <div class="text-caption text-grey-5 tw-mt-1">Try a different search term</div>
      </div>

      <div v-else-if="results.length === 0" class="empty-state tw-text-center tw-py-12">
        <div class="empty-icon tw-mx-auto tw-mb-4">
          <q-icon name="restaurant_menu" size="48px" color="grey-5" />
        </div>
        <div class="text-body1 text-grey-6">Start typing to search</div>
        <div class="text-caption text-grey-5 tw-mt-1">Find recipes by name or description</div>
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
        (r.description?.toLowerCase().includes(query.toLowerCase()) ?? false)
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
.search-page {
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 24px 16px;
  margin-bottom: 16px;
  color: white;
}

.header-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-input {
  :deep(.q-field__control) {
    border-radius: 14px;
  }

  :deep(.q-field__native) {
    color: #333;
  }
}

.empty-state {
  background: white;
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: #f5f5f5;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
