<template>
  <q-page class="mealplan-page">
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div class="tw-flex tw-items-center tw-gap-3">
          <div class="header-icon">
            <q-icon name="calendar_month" size="22px" color="white" />
          </div>
          <h1 class="page-title">Meal Plan</h1>
        </div>
        <div class="tw-flex tw-gap-2">
          <q-btn flat no-caps label="Today" class="header-btn" @click="goToCurrentWeek" />
          <q-btn
            flat
            round
            icon="shopping_cart"
            class="header-btn-icon"
            :loading="generatingShoppingList"
            @click="generateShoppingList"
          >
            <q-tooltip>Generate Shopping List</q-tooltip>
          </q-btn>
          <q-btn flat round icon="delete_sweep" class="header-btn-icon" @click="confirmClearWeek" />
        </div>
      </div>
    </header>

    <div class="tw-px-4 tw-pb-24">
      <WeekView
        :week-plan="currentWeek"
        @prev="navigateWeek('prev')"
        @next="navigateWeek('next')"
        @slot-click="openRecipeSelector"
        @slot-clear="clearSlot($event.id)"
      />
    </div>

    <!-- Recipe Selector Dialog -->
    <q-dialog v-model="selectorOpen">
      <div class="selector-dialog">
        <div class="selector-header">
          <div class="tw-flex tw-items-center tw-gap-3">
            <div class="selector-icon">
              <q-icon name="restaurant_menu" size="18px" color="white" />
            </div>
            <h2 class="selector-title">Select a Recipe</h2>
          </div>
          <q-btn flat round icon="close" size="sm" color="white" v-close-popup />
        </div>

        <div class="selector-content">
          <!-- Search & AI Suggest -->
          <div class="tw-flex tw-gap-3 tw-mb-4">
            <q-input
              v-model="searchQuery"
              dense
              outlined
              placeholder="Search recipes..."
              class="search-input tw-flex-1"
            >
              <template #prepend>
                <q-icon name="search" color="grey-6" />
              </template>
            </q-input>
            <q-btn
              icon="auto_awesome"
              unelevated
              no-caps
              class="suggest-btn"
              :loading="suggestionsLoading"
              @click="handleGetSuggestions"
            />
          </div>

          <!-- Error Banner -->
          <div v-if="suggestionsError" class="error-banner tw-mb-4">
            <q-icon name="error_outline" size="18px" />
            <span>{{ suggestionsError }}</span>
          </div>

          <!-- AI Suggestions -->
          <div v-if="suggestedRecipes.length > 0" class="suggestions-section tw-mb-4">
            <div class="tw-flex tw-items-center tw-justify-between tw-mb-3">
              <span class="section-label">
                <q-icon name="auto_awesome" size="16px" class="tw-mr-1" />
                AI Suggestions
              </span>
              <q-btn flat dense size="sm" icon="close" color="grey" @click="clearSuggestions" />
            </div>
            <div class="recipe-list suggestions-list">
              <div
                v-for="recipe in suggestedRecipes"
                :key="'suggestion-' + recipe.id"
                class="recipe-item recipe-item--suggested"
                @click="selectRecipe(recipe)"
              >
                <div class="recipe-icon">
                  <q-icon name="auto_awesome" size="16px" color="white" />
                </div>
                <div class="recipe-info">
                  <span class="recipe-name">{{ recipe.name }}</span>
                  <span class="recipe-desc">{{ recipe.description }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- All Recipes -->
          <div class="all-recipes-section">
            <span class="section-label">All Recipes</span>
            <div class="recipe-list">
              <div
                v-for="recipe in filteredRecipes"
                :key="recipe.id ?? ''"
                class="recipe-item"
                @click="selectRecipe(recipe)"
              >
                <div class="recipe-icon recipe-icon--default">
                  <q-icon name="restaurant_menu" size="16px" color="white" />
                </div>
                <div class="recipe-info">
                  <span class="recipe-name">{{ recipe.name }}</span>
                  <span class="recipe-desc">{{ recipe.description }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import WeekView from '@features/mealplan/components/WeekView.vue';
import { useMealplan } from '@features/mealplan/composables/useMealplan';
import { useRecipeStore } from '@features/recipe/store/recipeStore';
import { useCreateShoppingList } from '@features/shoppinglist/composables/useShoppinglist';
import type { HandlerRecipeJSON } from '@/api/generated/models';
import type { MealSlot } from '@features/mealplan/types/mealplan';

const $q = useQuasar();
const router = useRouter();
const recipeStore = useRecipeStore();
const { creating: generatingShoppingList, createFromRecipes } = useCreateShoppingList();
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

async function generateShoppingList() {
  // Get unique recipe IDs from the current week's meal plan
  const recipeIds = currentWeek.value.days
    .flatMap((day) => day.slots)
    .filter((slot) => slot.recipe?.id)
    .map((slot) => slot.recipe!.id!)
    .filter((id, index, self) => self.indexOf(id) === index); // unique

  if (recipeIds.length === 0) {
    $q.notify({
      type: 'warning',
      message: 'No recipes in your meal plan. Add some recipes first!',
    });
    return;
  }

  try {
    const list = await createFromRecipes({
      name: `Meal Plan - ${currentWeek.value.label}`,
      recipeIds,
    });
    if (list) {
      $q.notify({
        type: 'positive',
        message: `Shopping list created with ${list.totalItems} items`,
      });
      void router.push({ name: 'shoppinglist-detail', params: { id: list.id } });
    }
  } catch {
    $q.notify({
      type: 'negative',
      message: 'Failed to generate shopping list',
    });
  }
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600&display=swap');

.mealplan-page {
  padding-top: env(safe-area-inset-top);
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 20px 16px 24px;
  color: white;
}

.header-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-title {
  font-family: 'Fraunces', serif;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
  letter-spacing: -0.3px;
}

.header-btn {
  font-family: 'DM Sans', sans-serif;
  font-weight: 600;
  color: white;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 10px;

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }
}

.header-btn-icon {
  color: white;
  background: rgba(255, 255, 255, 0.15);

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }
}

// Selector Dialog
.selector-dialog {
  background: white;
  border-radius: 24px;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.selector-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: white;
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

.selector-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.selector-content {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.search-input {
  :deep(.q-field__control) {
    border-radius: 12px;
  }

  :deep(.q-field__native) {
    font-family: 'DM Sans', sans-serif;
  }
}

.suggest-btn {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  color: white !important;
  border-radius: 12px;
  width: 48px;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #ffebee;
  border-radius: 12px;
  color: #c62828;
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
}

.section-label {
  display: flex;
  align-items: center;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #ff6347;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
}

.recipe-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.suggestions-list {
  background: #fff5f2;
  border-radius: 16px;
  padding: 8px;
}

.all-recipes-section {
  .recipe-list {
    max-height: 300px;
    overflow-y: auto;
    background: #faf8f7;
    border-radius: 16px;
    padding: 8px;
  }
}

.recipe-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s ease;

  &:hover {
    background: rgba(255, 127, 80, 0.1);
  }

  &--suggested {
    background: white;
  }
}

.recipe-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &--default {
    background: linear-gradient(135deg, #ffa07a 0%, #ff7f50 100%);
  }
}

.recipe-info {
  flex: 1;
  min-width: 0;
}

.recipe-name {
  display: block;
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #2d1f1a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recipe-desc {
  display: block;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #6b5f5a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}
</style>
