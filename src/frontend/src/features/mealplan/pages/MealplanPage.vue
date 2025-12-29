<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
      <h1 class="text-h4 q-ma-none">Meal Plan</h1>
      <div class="row q-gutter-sm">
        <q-btn flat label="Today" @click="goToCurrentWeek" />
        <q-btn flat icon="delete_sweep" label="Clear Week" @click="confirmClearWeek" />
      </div>
    </div>

    <WeekView
      :week-plan="currentWeek"
      @prev="navigateWeek('prev')"
      @next="navigateWeek('next')"
      @slot-click="openRecipeSelector"
      @slot-clear="clearSlot($event.id)"
      @slot-suggest="openSuggestions"
    />

    <q-dialog v-model="selectorOpen">
      <q-card style="min-width: 400px; max-width: 600px">
        <q-card-section>
          <div class="text-h6">Select a Recipe</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input v-model="searchQuery" dense outlined placeholder="Search recipes..." />
        </q-card-section>

        <q-card-section class="scroll" style="max-height: 400px">
          <q-list separator>
            <q-item
              v-for="recipe in filteredRecipes"
              :key="recipe.id"
              clickable
              v-ripple
              @click="selectRecipe(recipe)"
            >
              <q-item-section>
                <q-item-label>{{ recipe.name }}</q-item-label>
                <q-item-label caption>{{ recipe.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <SuggestionsDialog
      v-model="suggestionsOpen"
      :meal-slot="suggestingSlot"
      :suggestions="suggestions"
      :loading="suggestionsLoading"
      :error="error"
      @select="selectSuggestion"
      @retry="retrySuggestions"
      @update:model-value="!$event && closeSuggestions()"
    />
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { WeekView, SuggestionsDialog } from '../components';
import { useMealplan } from '../composables';
import { useRecipeStore } from '@/features/recipe/store';
import type { Recipe } from '@/features/recipe/types';
import type { MealSlot } from '../types';

const $q = useQuasar();
const recipeStore = useRecipeStore();
const {
  currentWeek,
  setRecipeForSlot,
  clearSlot,
  navigateWeek,
  goToCurrentWeek,
  clearWeek,
  suggestions,
  suggestionsLoading,
  error,
  fetchSuggestions,
  clearSuggestions,
} = useMealplan();

const selectorOpen = ref(false);
const selectedSlot = ref<MealSlot | null>(null);
const searchQuery = ref('');
const suggestionsOpen = ref(false);
const suggestingSlot = ref<MealSlot | null>(null);

const filteredRecipes = computed(() => {
  const query = searchQuery.value.toLowerCase();
  if (!query) return recipeStore.recipes;
  return recipeStore.recipes.filter(
    (r) =>
      r.name.toLowerCase().includes(query) ||
      r.description.toLowerCase().includes(query)
  );
});

onMounted(() => {
  if (recipeStore.recipes.length === 0) {
    recipeStore.fetchRecipes();
  }
});

function openRecipeSelector(slot: MealSlot) {
  selectedSlot.value = slot;
  searchQuery.value = '';
  selectorOpen.value = true;
}

function selectRecipe(recipe: Recipe) {
  if (selectedSlot.value) {
    setRecipeForSlot(selectedSlot.value.id, recipe);
  }
  selectorOpen.value = false;
  selectedSlot.value = null;
}

function confirmClearWeek() {
  $q.dialog({
    title: 'Clear Week',
    message: 'Are you sure you want to clear all meals for this week?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    clearWeek();
  });
}

function openSuggestions(slot: MealSlot) {
  suggestingSlot.value = slot;
  suggestionsOpen.value = true;
  void fetchSuggestions();
}

function closeSuggestions() {
  suggestionsOpen.value = false;
  suggestingSlot.value = null;
  clearSuggestions();
}

function selectSuggestion(recipe: Recipe) {
  if (suggestingSlot.value) {
    setRecipeForSlot(suggestingSlot.value.id, recipe);
  }
  closeSuggestions();
}

function retrySuggestions() {
  void fetchSuggestions();
}
</script>
