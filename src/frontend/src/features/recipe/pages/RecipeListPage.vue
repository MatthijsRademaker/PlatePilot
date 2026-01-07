<template>
  <q-page class="recipe-list-page">
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div class="tw-flex tw-items-center tw-gap-3">
          <div class="header-icon">
            <q-icon name="menu_book" size="22px" color="white" />
          </div>
          <h1 class="page-title">Recipes</h1>
        </div>
        <q-btn
          icon="refresh"
          flat
          round
          class="refresh-btn"
          aria-label="Refresh"
          @click="refresh"
        />
      </div>
    </header>

    <div class="tw-px-4 tw-pb-24">
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
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600;9..144,700&display=swap');

.recipe-list-page {
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

.refresh-btn {
  color: white;
  background: rgba(255, 255, 255, 0.15);

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }
}
</style>
