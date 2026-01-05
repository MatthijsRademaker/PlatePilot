<template>
  <q-page class="mealplan-page">
    <div class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div class="tw-flex tw-items-center tw-gap-3">
          <div class="header-icon">
            <q-icon name="calendar_month" size="24px" color="white" />
          </div>
          <h1 class="text-h5 q-ma-none tw-font-semibold">Meal Plan</h1>
        </div>
        <div class="row q-gutter-sm">
          <q-btn flat label="Today" class="header-btn" @click="goToCurrentWeek" />
          <q-btn flat icon="delete_sweep" class="header-btn" @click="confirmClearWeek" />
        </div>
      </div>
    </div>

    <div class="tw-px-4 tw-pb-4">
      <WeekView
        :week-plan="currentWeek"
        @prev="navigateWeek('prev')"
        @next="navigateWeek('next')"
        @slot-click="openRecipeSelector"
        @slot-clear="clearSlot($event.id)"
      />
    </div>

    <q-dialog v-model="selectorOpen">
      <q-card class="selector-card">
        <q-card-section class="selector-header">
          <div class="tw-flex tw-items-center tw-gap-2">
            <div class="selector-icon">
              <q-icon name="restaurant_menu" size="20px" color="white" />
            </div>
            <div class="text-h6 tw-font-semibold">Select a Recipe</div>
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="row q-gutter-sm q-mb-md">
            <q-input
              v-model="searchQuery"
              dense
              outlined
              placeholder="Search recipes..."
              class="col search-input"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
            <q-btn
              color="primary"
              icon="auto_awesome"
              label="AI Suggest"
              unelevated
              :loading="suggestionsLoading"
              @click="handleGetSuggestions"
            />
          </div>
        </q-card-section>

        <!-- AI Suggestions Section -->
        <q-card-section v-if="suggestionsError" class="q-pt-none">
          <q-banner class="bg-negative text-white tw-rounded-lg">
            {{ suggestionsError }}
          </q-banner>
        </q-card-section>

        <q-card-section v-if="suggestedRecipes.length > 0" class="q-pt-none">
          <div class="q-mb-sm row items-center justify-between">
            <span class="text-subtitle2 text-primary tw-flex tw-items-center tw-gap-1">
              <q-icon name="auto_awesome" />
              AI Suggestions
            </span>
            <q-btn flat dense size="sm" icon="close" @click="clearSuggestions" />
          </div>
          <q-list bordered separator class="suggestions-list tw-rounded-lg">
            <q-item
              v-for="recipe in suggestedRecipes"
              :key="'suggestion-' + recipe.id"
              clickable
              v-ripple
              @click="selectRecipe(recipe)"
            >
              <q-item-section avatar>
                <div class="suggestion-avatar">
                  <q-icon name="auto_awesome" size="16px" color="white" />
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="tw-font-medium">{{ recipe.name }}</q-item-label>
                <q-item-label caption>{{ recipe.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <!-- All Recipes Section -->
        <q-card-section class="scroll q-pt-none" style="max-height: 300px">
          <div class="text-subtitle2 text-grey-7 q-mb-sm">All Recipes</div>
          <q-list separator class="recipes-list tw-rounded-lg">
            <q-item
              v-for="recipe in filteredRecipes"
              :key="recipe.id ?? ''"
              clickable
              v-ripple
              @click="selectRecipe(recipe)"
            >
              <q-item-section avatar>
                <div class="recipe-avatar">
                  <q-icon name="restaurant_menu" size="16px" color="white" />
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="tw-font-medium">{{ recipe.name }}</q-item-label>
                <q-item-label caption>{{ recipe.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="tw-px-4 tw-pb-4">
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
import type { HandlerRecipeJSON } from '@/api/generated/models';
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
  suggestedRecipeIds,
  suggestionsLoading,
  suggestionsError,
  fetchSuggestions,
  clearSuggestions,
} = useMealplan();

const selectorOpen = ref(false);
const selectedSlot = ref<MealSlot | null>(null);
const searchQuery = ref('');

// Compute suggested recipes from IDs by looking them up in the loaded recipes
const suggestedRecipes = computed(() => {
  const ids = suggestedRecipeIds.value;
  if (!ids || ids.length === 0) return [];
  return recipeStore.recipes.filter((r) => r.id && ids.includes(r.id));
});

const filteredRecipes = computed(() => {
  const query = searchQuery.value.toLowerCase();
  if (!query) return recipeStore.recipes;
  return recipeStore.recipes.filter(
    (r) =>
      (r.name?.toLowerCase().includes(query) ?? false) ||
      (r.description?.toLowerCase().includes(query) ?? false)
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

function selectRecipe(recipe: HandlerRecipeJSON) {
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

<style scoped lang="scss">
.mealplan-page {
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 24px 16px;
  margin-bottom: 16px;
  color: white;
}

.header-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-btn {
  color: white;
  background: rgba(255, 255, 255, 0.15);

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }
}

.selector-card {
  min-width: 400px;
  max-width: 600px;
  border-radius: 20px;
}

.selector-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  color: white;
  border-radius: 20px 20px 0 0;
}

.selector-icon {
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-input {
  :deep(.q-field__control) {
    border-radius: 12px;
  }
}

.suggestions-list {
  background: #fff5f2;
  border-color: rgba(255, 127, 80, 0.2) !important;
}

.suggestion-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recipes-list {
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.recipe-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #ffa07a 0%, #ff7f50 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
