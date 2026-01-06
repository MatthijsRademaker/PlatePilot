<template>
  <q-page class="recipe-list-page">
    <div class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div class="tw-flex tw-items-center tw-gap-3">
          <div class="header-icon">
            <q-icon name="menu_book" size="24px" color="white" />
          </div>
          <h1 class="text-h5 q-ma-none tw-font-semibold">Recipes</h1>
        </div>
        <q-btn color="primary" icon="refresh" flat round aria-label="Refresh" @click="refresh" />
      </div>
    </div>

    <div class="tw-px-4 tw-pb-4 recipe-list">
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
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useRecipeList } from '@features/recipe/composables/useRecipe';
import RecipeList from '@features/recipe/components/RecipeList.vue';
import type { HandlerRecipeJSON } from '@/api/generated/models';

const router = useRouter();
const { recipes, loading, error, pageIndex, totalPages, loadPage, refresh } = useRecipeList();

function goToRecipe(recipe: HandlerRecipeJSON) {
  if (recipe.id) {
    void router.push({ name: 'recipe-detail', params: { id: recipe.id } });
  }
}
</script>

<style scoped lang="scss">
.recipe-list-page {
  padding-top: env(safe-area-inset-top);
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.recipe-list {
  border-top-right-radius: 8px;
  border-top-left-radius: 8px;
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
</style>
