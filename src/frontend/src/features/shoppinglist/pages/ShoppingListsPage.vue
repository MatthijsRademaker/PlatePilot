<template>
  <q-page class="shopping-lists-page">
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-justify-between">
        <div class="tw-flex tw-items-center tw-gap-3">
          <div class="header-icon">
            <q-icon name="shopping_cart" size="22px" color="white" />
          </div>
          <h1 class="page-title">Shopping Lists</h1>
        </div>
        <q-btn
          icon="add"
          flat
          round
          class="add-btn"
          aria-label="Create List"
          @click="showCreateDialog = true"
        />
      </div>
    </header>

    <div class="tw-px-4 tw-pb-24 tw-pt-4">
      <!-- Loading State -->
      <div v-if="loading" class="tw-text-center tw-py-12">
        <q-spinner size="40px" color="primary" />
        <p class="tw-mt-4 tw-text-gray-500">Loading lists...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="tw-text-center tw-py-12">
        <q-icon name="error_outline" size="48px" color="negative" />
        <p class="tw-mt-4 tw-text-red-600">{{ error }}</p>
        <q-btn
          label="Retry"
          color="primary"
          outline
          class="tw-mt-4"
          @click="refresh"
        />
      </div>

      <!-- Empty State -->
      <div v-else-if="lists.length === 0" class="empty-state">
        <div class="empty-icon">
          <q-icon name="shopping_cart" size="48px" color="grey-5" />
        </div>
        <h3 class="empty-title">No shopping lists yet</h3>
        <p class="empty-description">
          Create a shopping list from your meal plan or start a new list manually.
        </p>
        <q-btn
          label="Create List"
          color="primary"
          unelevated
          class="tw-mt-4"
          @click="showCreateDialog = true"
        />
      </div>

      <!-- List Grid -->
      <div v-else class="lists-grid">
        <ShoppingListCard
          v-for="list in lists"
          :key="list.id"
          :shopping-list="list"
          @click="goToList"
        />
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="tw-flex tw-justify-center tw-mt-6">
        <q-pagination
          :model-value="pageIndex"
          :max="totalPages"
          direction-links
          flat
          color="primary"
          @update:model-value="loadPage"
        />
      </div>
    </div>

    <!-- Create Dialog -->
    <q-dialog v-model="showCreateDialog">
      <q-card style="min-width: 350px; border-radius: 16px">
        <q-card-section>
          <div class="dialog-title">Create Shopping List</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="newListName"
            label="List Name"
            outlined
            dense
            autofocus
            placeholder="e.g., Weekly Groceries"
            @keyup.enter="handleCreate"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-4 q-pb-4">
          <q-btn flat label="Cancel" color="grey" @click="showCreateDialog = false" />
          <q-btn
            unelevated
            label="Create"
            color="primary"
            :loading="creating"
            @click="handleCreate"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useShoppinglistList } from '../composables/useShoppinglist';
import ShoppingListCard from '../components/ShoppingListCard.vue';
import type { ShoppingListSummaryJSON } from '../types/shoppinglist';

const router = useRouter();
const { lists, loading, error, pageIndex, totalPages, loadPage, refresh, createEmptyList } =
  useShoppinglistList();

const showCreateDialog = ref(false);
const newListName = ref('');
const creating = ref(false);

function goToList(list: ShoppingListSummaryJSON) {
  void router.push({ name: 'shoppinglist-detail', params: { id: list.id } });
}

async function handleCreate() {
  creating.value = true;
  try {
    const newList = await createEmptyList(newListName.value || undefined);
    showCreateDialog.value = false;
    newListName.value = '';
    if (newList) {
      void router.push({ name: 'shoppinglist-detail', params: { id: newList.id } });
    }
  } catch (e) {
    // Error handled in store
  } finally {
    creating.value = false;
  }
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600;9..144,700&display=swap');

.shopping-lists-page {
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

.add-btn {
  color: white;
  background: rgba(255, 255, 255, 0.15);

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }
}

.lists-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: #f0ece9;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.empty-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0 0 8px;
}

.empty-description {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  color: #6b5f5a;
  margin: 0;
  max-width: 280px;
  margin: 0 auto;
}

.dialog-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
}
</style>
