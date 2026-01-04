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
    />

    <q-dialog v-model="selectorOpen">
      <q-card style="min-width: 400px; max-width: 600px">
        <q-card-section>
          <div class="text-h6">Select a Recipe</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="row q-gutter-sm q-mb-md">
            <q-input
              v-model="searchQuery"
              dense
              outlined
              placeholder="Search recipes..."
              class="col"
            />
            <q-btn
              color="primary"
              icon="auto_awesome"
              label="Get AI Suggestions"
              :loading="suggestionsLoading"
              @click="handleGetSuggestions"
            />
          </div>
        </q-card-section>

        <!-- AI Suggestions Section -->
        <q-card-section v-if="suggestionsError" class="q-pt-none">
          <q-banner class="bg-negative text-white">
            {{ suggestionsError }}
          </q-banner>
        </q-card-section>

        <q-card-section v-if="suggestions.length > 0" class="q-pt-none">
          <div class="q-mb-sm row items-center justify-between">
            <span class="text-subtitle2 text-primary">
              <q-icon name="auto_awesome" class="q-mr-xs" />
              AI Suggestions
            </span>
            <q-btn flat dense size="sm" icon="close" @click="clearSuggestions" />
          </div>
          <q-list bordered separator class="rounded-borders bg-blue-1">
            <q-item
              v-for="recipe in suggestions"
              :key="'suggestion-' + recipe.id"
              clickable
              v-ripple
              @click="selectRecipe(recipe)"
            >
              <q-item-section avatar>
                <q-icon name="auto_awesome" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ recipe.name }}</q-item-label>
                <q-item-label caption>{{ recipe.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <!-- All Recipes Section -->
        <q-card-section class="scroll q-pt-none" style="max-height: 300px">
          <div class="text-subtitle2 text-grey-7 q-mb-sm">All Recipes</div>
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
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import WeekView from '@features/mealplan/components/WeekView.vue';
import { useMealplan } from '@features/mealplan/composables/useMealplan';
import { useRecipeStore } from '@features/recipe/store/recipeStore';
import type { Recipe } from '@features/recipe/types/recipe';
import type { MealSlot } from '@features/mealplan/types/mealplan';

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
  suggestionsError,
  fetchSuggestions,
  clearSuggestions,
} = useMealplan();

const selectorOpen = ref(false);
const selectedSlot = ref<MealSlot | null>(null);
const searchQuery = ref('');

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
    void recipeStore.fetchRecipes();
  }
});

function openRecipeSelector(slot: MealSlot) {
  selectedSlot.value = slot;
  searchQuery.value = '';
  clearSuggestions();
  selectorOpen.value = true;
}

function selectRecipe(recipe: Recipe) {
  if (selectedSlot.value) {
    setRecipeForSlot(selectedSlot.value.id, recipe);
  }
  selectorOpen.value = false;
  selectedSlot.value = null;
}

function handleGetSuggestions() {
  void fetchSuggestions(5);
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
</script>
