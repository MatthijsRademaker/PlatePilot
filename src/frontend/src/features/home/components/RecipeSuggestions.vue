<template>
  <div class="recipe-suggestions">
    <div class="tw-text-lg tw-font-semibold tw-text-gray-800 tw-mb-4">Recipe Suggestions</div>

    <div v-if="loading" class="tw-flex tw-gap-3 tw-overflow-x-auto tw-pb-2">
      <div v-for="n in 4" :key="n" class="suggestion-card-skeleton">
        <q-skeleton height="100px" class="tw-rounded-t-xl" />
        <div class="tw-p-2">
          <q-skeleton type="text" width="80%" />
        </div>
      </div>
    </div>

    <div v-else-if="recipes.length > 0" class="suggestions-scroll tw-flex tw-gap-3 tw-overflow-x-auto tw-pb-2">
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        class="suggestion-card"
        @click="viewRecipe(recipe.id)"
      >
        <div class="card-image">
          <img :src="getRecipeImage(recipe.name)" :alt="recipe.name" />
        </div>
        <div class="tw-p-2">
          <div class="tw-text-sm tw-font-medium tw-text-gray-800 tw-truncate">
            {{ recipe.name }}
          </div>
          <div v-if="recipe.cuisine?.name" class="tw-text-xs tw-text-gray-500">
            {{ recipe.cuisine.name }}
          </div>
        </div>
      </div>
    </div>

    <div v-else class="tw-text-center tw-py-6 tw-text-gray-500">
      No recipes available
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import type { HandlerRecipeJSON } from '@/api/generated/models';
import { useRecipeStore } from '@features/recipe/store/recipeStore';

const router = useRouter();
const recipeStore = useRecipeStore();

const recipes = ref<HandlerRecipeJSON[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    // Fetch first page of recipes as suggestions
    await recipeStore.fetchRecipes(1);
    // Take only first 6 for suggestions
    recipes.value = recipeStore.recipes.slice(0, 6);
  } finally {
    loading.value = false;
  }
});

function getRecipeImage(recipeName: string | undefined): string {
  const seed = recipeName?.replace(/\s+/g, '-').toLowerCase() || 'default';
  return `https://picsum.photos/seed/${seed}/200/150`;
}

function viewRecipe(id: string | undefined) {
  if (id) {
    void router.push({ name: 'recipe-detail', params: { id } });
  }
}
</script>

<style scoped lang="scss">
.recipe-suggestions {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.suggestions-scroll {
  margin: 0 -16px;
  padding: 0 16px;

  &::-webkit-scrollbar {
    display: none;
  }

  -ms-overflow-style: none;
  scrollbar-width: none;
}

.suggestion-card {
  flex-shrink: 0;
  width: 140px;
  background: #fafafa;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 127, 80, 0.15);
  }
}

.suggestion-card-skeleton {
  flex-shrink: 0;
  width: 140px;
  background: #fafafa;
  border-radius: 12px;
  overflow: hidden;
}

.card-image {
  width: 100%;
  height: 100px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
