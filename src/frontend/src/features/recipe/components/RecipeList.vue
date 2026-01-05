<template>
  <div class="recipe-list">
    <div v-if="loading && recipes.length === 0" class="row justify-center tw-py-8">
      <q-spinner size="lg" color="primary" />
    </div>

    <div v-else-if="error" class="empty-state tw-text-center tw-py-12">
      <div class="empty-icon empty-icon--error tw-mx-auto tw-mb-4">
        <q-icon name="error" size="48px" color="negative" />
      </div>
      <p class="text-negative tw-mb-4">{{ error }}</p>
      <q-btn color="primary" label="Retry" unelevated rounded @click="$emit('refresh')" />
    </div>

    <div v-else-if="recipes.length === 0" class="empty-state tw-text-center tw-py-12">
      <div class="empty-icon tw-mx-auto tw-mb-4">
        <q-icon name="restaurant_menu" size="48px" color="grey-5" />
      </div>
      <div class="text-body1 text-grey-6">No recipes found</div>
      <div class="text-caption text-grey-5 tw-mt-1">Check back later for new recipes</div>
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
        color="primary"
        active-color="primary"
        @update:model-value="$emit('page', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { HandlerRecipeJSON } from '@/api/generated/models';
import RecipeCard from '@features/recipe/components/RecipeCard.vue';

defineProps<{
  recipes: HandlerRecipeJSON[];
  loading: boolean;
  error: string | null;
  pageIndex: number;
  totalPages: number;
}>();

defineEmits<{
  select: [recipe: HandlerRecipeJSON];
  page: [page: number];
  refresh: [];
}>();
</script>

<style scoped lang="scss">
.empty-state {
  background: white;
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: #f5f5f5;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;

  &--error {
    background: #ffebee;
  }
}
</style>
