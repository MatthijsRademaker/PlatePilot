<template>
  <q-card class="recipe-card" @click="$emit('click', recipe)">
    <q-img v-if="recipe.imageUrl" :src="recipe.imageUrl" :ratio="16 / 9" />
    <div v-else class="placeholder-image">
      <q-icon name="restaurant_menu" size="64px" color="grey-5" />
    </div>

    <q-card-section>
      <div class="text-h6 ellipsis">{{ recipe.name }}</div>
      <div class="text-caption text-grey ellipsis-2-lines">
        {{ recipe.description }}
      </div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <div class="row items-center q-gutter-sm">
        <q-chip dense size="sm" icon="schedule">
          {{ totalTime }} min
        </q-chip>
        <q-chip dense size="sm" icon="restaurant">
          {{ recipe.servings }} servings
        </q-chip>
      </div>
    </q-card-section>

    <q-card-section v-if="recipe.cuisines.length > 0" class="q-pt-none">
      <q-chip
        v-for="cuisine in recipe.cuisines.slice(0, 3)"
        :key="cuisine.id"
        dense
        size="sm"
        color="primary"
        text-color="white"
      >
        {{ cuisine.name }}
      </q-chip>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Recipe } from '../types';

const props = defineProps<{
  recipe: Recipe;
}>();

defineEmits<{
  click: [recipe: Recipe];
}>();

const totalTime = computed(() => props.recipe.preparationTime + props.recipe.cookingTime);
</script>

<style scoped lang="scss">
.recipe-card {
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-4px);
  }
}

.placeholder-image {
  aspect-ratio: 16 / 9;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--q-grey-3);
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
