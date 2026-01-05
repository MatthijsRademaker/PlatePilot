<template>
  <q-card class="recipe-card" @click="$emit('click', recipe)">
    <div class="placeholder-image">
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
        <q-chip v-if="recipe.prepTime" dense size="sm" icon="schedule">
          Prep: {{ recipe.prepTime }}
        </q-chip>
        <q-chip v-if="recipe.cookTime" dense size="sm" icon="local_fire_department">
          Cook: {{ recipe.cookTime }}
        </q-chip>
      </div>
    </q-card-section>

    <q-card-section v-if="recipe.cuisine" class="q-pt-none">
      <q-chip dense size="sm" color="primary" text-color="white">
        {{ recipe.cuisine.name }}
      </q-chip>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import type { HandlerRecipeJSON } from '@/api/generated/models';

defineProps<{
  recipe: HandlerRecipeJSON;
}>();

defineEmits<{
  click: [recipe: HandlerRecipeJSON];
}>();
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
