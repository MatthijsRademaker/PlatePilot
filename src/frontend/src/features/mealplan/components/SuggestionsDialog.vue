<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card class="suggestions-dialog">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">
          <q-icon name="auto_awesome" color="primary" class="q-mr-sm" />
          Recipe Suggestions
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section v-if="mealSlot" class="q-pt-sm">
        <q-chip dense icon="event" color="grey-3">
          {{ formatSlotLabel(mealSlot) }}
        </q-chip>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <!-- Loading state -->
        <div v-if="loading" class="column items-center q-pa-lg">
          <q-spinner-dots color="primary" size="40px" />
          <div class="text-body2 text-grey q-mt-md">Finding recipes for you...</div>
        </div>

        <!-- Error state -->
        <div v-else-if="error" class="column items-center q-pa-lg">
          <q-icon name="error_outline" color="negative" size="48px" />
          <div class="text-body2 text-negative q-mt-md">{{ error }}</div>
          <q-btn
            flat
            color="primary"
            label="Try Again"
            class="q-mt-md"
            @click="$emit('retry')"
          />
        </div>

        <!-- Empty state -->
        <div v-else-if="suggestions.length === 0" class="column items-center q-pa-lg">
          <q-icon name="restaurant_menu" color="grey" size="48px" />
          <div class="text-body2 text-grey q-mt-md">No suggestions available</div>
        </div>

        <!-- Suggestions list -->
        <q-list v-else separator class="suggestions-list">
          <q-item
            v-for="recipe in suggestions"
            :key="recipe.id"
            clickable
            v-ripple
            class="suggestion-item"
            @click="$emit('select', recipe)"
          >
            <q-item-section avatar v-if="recipe.imageUrl">
              <q-avatar rounded size="56px">
                <q-img :src="recipe.imageUrl" />
              </q-avatar>
            </q-item-section>
            <q-item-section avatar v-else>
              <q-avatar rounded size="56px" color="grey-3" text-color="grey">
                <q-icon name="restaurant" />
              </q-avatar>
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-medium">{{ recipe.name }}</q-item-label>
              <q-item-label caption lines="2">{{ recipe.description }}</q-item-label>
              <q-item-label caption class="q-mt-xs">
                <q-icon name="schedule" size="xs" class="q-mr-xs" />
                {{ recipe.preparationTime + recipe.cookingTime }} min
                <span class="q-mx-sm">|</span>
                <q-icon name="restaurant" size="xs" class="q-mr-xs" />
                {{ recipe.servings }} servings
              </q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-icon name="add_circle" color="primary" size="24px" />
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <q-card-actions v-if="!loading && suggestions.length > 0" align="right" class="q-pt-none">
        <q-btn flat label="Cancel" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import type { Recipe } from '@/features/recipe/types';
import type { MealSlot } from '../types';

defineProps<{
  modelValue: boolean;
  mealSlot: MealSlot | null;
  suggestions: Recipe[];
  loading: boolean;
  error: string | null;
}>();

defineEmits<{
  'update:modelValue': [value: boolean];
  select: [recipe: Recipe];
  retry: [];
}>();

function formatSlotLabel(slot: MealSlot): string {
  const date = new Date(slot.date);
  const dayName = date.toLocaleDateString('en-US', { weekday: 'long' });
  const mealType = slot.mealType.charAt(0).toUpperCase() + slot.mealType.slice(1);
  return `${dayName} - ${mealType}`;
}
</script>

<style scoped lang="scss">
.suggestions-dialog {
  min-width: 400px;
  max-width: 500px;
  width: 100%;
}

.suggestions-list {
  max-height: 400px;
  overflow-y: auto;
}

.suggestion-item {
  border-radius: 8px;
  margin: 4px 0;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(0, 0, 0, 0.04);
  }
}
</style>
