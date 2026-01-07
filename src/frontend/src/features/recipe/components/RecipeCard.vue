<template>
  <div class="recipe-card" @click="$emit('click', recipe)">
    <div class="card-image">
      <img :src="getRecipeImage(recipe.name)" :alt="recipe.name" />
      <div v-if="recipe.cuisine" class="cuisine-badge">
        {{ recipe.cuisine.name }}
      </div>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ recipe.name }}</h3>
      <p v-if="recipe.description" class="card-description">
        {{ recipe.description }}
      </p>

      <div v-if="recipe.prepTime || recipe.cookTime" class="card-meta">
        <span v-if="recipe.prepTime" class="meta-item">
          <q-icon name="schedule" size="14px" />
          {{ recipe.prepTime }}
        </span>
        <span v-if="recipe.cookTime" class="meta-item">
          <q-icon name="local_fire_department" size="14px" />
          {{ recipe.cookTime }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { HandlerRecipeJSON } from '@/api/generated/models';

defineProps<{
  recipe: HandlerRecipeJSON;
}>();

defineEmits<{
  click: [recipe: HandlerRecipeJSON];
}>();

function getRecipeImage(recipeName: string | undefined): string {
  const seed = recipeName?.replace(/\s+/g, '-').toLowerCase() || 'default';
  return `https://picsum.photos/seed/${seed}/400/300`;
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600&display=swap');

.recipe-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
  border: 1px solid rgba(45, 31, 26, 0.04);
  box-shadow: 0 2px 12px rgba(45, 31, 26, 0.04);

  &:active {
    transform: scale(0.98);
  }

  @media (hover: hover) {
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 32px rgba(255, 127, 80, 0.15);
    }

    &:hover .card-image img {
      transform: scale(1.05);
    }
  }
}

.card-image {
  position: relative;
  aspect-ratio: 16 / 10;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
  }
}

.cuisine-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 6px 12px;
  border-radius: 20px;
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 600;
  color: white;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.card-content {
  padding: 16px;
}

.card-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0 0 6px;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-description {
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 400;
  color: #6b5f5a;
  margin: 0 0 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  gap: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 500;
  color: #ff6347;
  background: #fff5f2;
  padding: 4px 10px;
  border-radius: 8px;
}
</style>
