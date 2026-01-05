<template>
  <q-card class="recipe-card" flat @click="$emit('click', recipe)">
    <div class="placeholder-image">
      <q-icon name="restaurant_menu" size="48px" color="white" />
    </div>

    <q-card-section>
      <div class="text-h6 ellipsis tw-font-semibold">{{ recipe.name }}</div>
      <div class="text-caption text-grey-7 ellipsis-2-lines tw-mt-1">
        {{ recipe.description }}
      </div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <div class="row items-center q-gutter-sm">
        <q-chip
          v-if="recipe.prepTime"
          dense
          size="sm"
          icon="schedule"
          class="time-chip"
        >
          {{ recipe.prepTime }}
        </q-chip>
        <q-chip
          v-if="recipe.cookTime"
          dense
          size="sm"
          icon="local_fire_department"
          class="time-chip"
        >
          {{ recipe.cookTime }}
        </q-chip>
      </div>
    </q-card-section>

    <q-card-section v-if="recipe.cuisine" class="q-pt-none">
      <q-chip dense size="sm" class="cuisine-chip">
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
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  overflow: hidden;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(255, 127, 80, 0.15);
  }
}

.placeholder-image {
  aspect-ratio: 16 / 9;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ffa07a 0%, #ff7f50 100%);
}

.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.time-chip {
  background: #fff5f2 !important;
  color: #ff6347 !important;
}

.cuisine-chip {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
}
</style>
