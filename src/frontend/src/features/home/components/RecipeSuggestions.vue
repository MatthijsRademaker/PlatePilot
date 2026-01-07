<template>
  <div class="recipe-suggestions">
    <div class="section-header">
      <h2 class="section-title">Recipe Suggestions</h2>
      <router-link :to="{ name: 'recipes' }" class="see-all">
        See all
        <q-icon name="chevron_right" size="18px" />
      </router-link>
    </div>

    <div v-if="loading" class="suggestions-scroll">
      <div v-for="n in 4" :key="n" class="suggestion-card suggestion-card--skeleton">
        <q-skeleton height="100px" class="skeleton-image" />
        <div class="skeleton-content">
          <q-skeleton type="text" width="80%" />
          <q-skeleton type="text" width="50%" class="tw-mt-1" />
        </div>
      </div>
    </div>

    <div v-else-if="recipes.length > 0" class="suggestions-scroll">
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        class="suggestion-card"
        @click="viewRecipe(recipe.id)"
      >
        <div class="card-image">
          <img :src="getRecipeImage(recipe.name)" :alt="recipe.name" />
        </div>
        <div class="card-content">
          <h3 class="card-title">{{ recipe.name }}</h3>
          <span v-if="recipe.cuisine?.name" class="card-cuisine">
            {{ recipe.cuisine.name }}
          </span>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <q-icon name="restaurant_menu" size="32px" color="grey-4" />
      <p>No recipes available</p>
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
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600&display=swap');

.recipe-suggestions {
  background: white;
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 31, 26, 0.06);
  border: 1px solid rgba(45, 31, 26, 0.04);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
}

.see-all {
  display: flex;
  align-items: center;
  gap: 2px;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 600;
  color: #ff6347;
  text-decoration: none;
  transition: opacity 0.2s ease;

  &:active {
    opacity: 0.7;
  }
}

.suggestions-scroll {
  display: flex;
  gap: 12px;
  overflow-x: auto;
  margin: 0 -20px;
  padding: 4px 20px 8px;

  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.suggestion-card {
  flex-shrink: 0;
  width: 140px;
  background: #faf8f7;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:active {
    transform: scale(0.98);
  }

  @media (hover: hover) {
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(255, 127, 80, 0.15);
    }
  }

  &--skeleton {
    pointer-events: none;
  }
}

.card-image {
  width: 100%;
  height: 100px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }
}

.suggestion-card:hover .card-image img {
  transform: scale(1.05);
}

.skeleton-image {
  border-radius: 0;
}

.skeleton-content {
  padding: 10px 12px;
}

.card-content {
  padding: 10px 12px;
}

.card-title {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-cuisine {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 500;
  color: #a8a0a0;
  display: block;
  margin-top: 2px;
}

.empty-state {
  text-align: center;
  padding: 32px 16px;

  p {
    font-family: 'DM Sans', sans-serif;
    font-size: 14px;
    color: #a8a0a0;
    margin: 12px 0 0;
  }
}
</style>
