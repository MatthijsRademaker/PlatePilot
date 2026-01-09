<template>
  <q-page class="shopping-list-detail-page">
    <!-- Header -->
    <header class="page-header">
      <div class="tw-flex tw-items-center tw-gap-3">
        <q-btn
          icon="arrow_back"
          flat
          round
          class="back-btn"
          @click="router.back()"
        />
        <div v-if="list" class="tw-flex-1 tw-min-w-0">
          <h1 class="page-title">{{ list.name }}</h1>
          <div class="progress-info">
            {{ list.checkedItems }} of {{ list.totalItems }} items
          </div>
        </div>
        <q-btn
          v-if="list"
          icon="more_vert"
          flat
          round
          class="menu-btn"
        >
          <q-menu>
            <q-list dense>
              <q-item clickable v-close-popup @click="showRenameDialog = true">
                <q-item-section avatar>
                  <q-icon name="edit" size="20px" />
                </q-item-section>
                <q-item-section>Rename</q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable v-close-popup @click="confirmDelete">
                <q-item-section avatar>
                  <q-icon name="delete" size="20px" color="negative" />
                </q-item-section>
                <q-item-section class="text-negative">Delete List</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>

      <!-- Progress Bar -->
      <div v-if="list" class="progress-bar-container">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
      </div>
    </header>

    <div class="tw-px-4 tw-pb-24 tw-pt-4">
      <!-- Loading State -->
      <div v-if="loading" class="tw-text-center tw-py-12">
        <q-spinner size="40px" color="primary" />
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
      <div v-else-if="list && list.items.length === 0" class="empty-state">
        <q-icon name="shopping_basket" size="48px" color="grey-5" />
        <p class="tw-mt-4 tw-text-gray-500">No items in this list</p>
        <q-btn
          label="Add Item"
          color="primary"
          unelevated
          class="tw-mt-4"
          @click="showAddDialog = true"
        />
      </div>

      <!-- Items List -->
      <div v-else-if="list">
        <!-- View Toggle -->
        <div class="view-toggle tw-mb-4">
          <q-btn-toggle
            v-model="groupByCategory"
            spread
            no-caps
            rounded
            toggle-color="primary"
            :options="[
              { label: 'Grouped', value: true },
              { label: 'All Items', value: false },
            ]"
          />
        </div>

        <!-- Grouped View -->
        <div v-if="groupByCategory">
          <ShoppingListGroup
            v-for="group in groupedItems"
            :key="group.categoryId"
            :group="group"
            :show-header="true"
            :show-item-actions="true"
            @toggle-item="handleToggle"
            @delete-item="handleDeleteItem"
            @view-sources="showSourcesFor"
          />
        </div>

        <!-- Flat View -->
        <div v-else class="flat-list">
          <ShoppingListItem
            v-for="item in list.items"
            :key="item.id"
            :item="item"
            :show-actions="true"
            @toggle="handleToggle"
            @delete="handleDeleteItem"
            @view-sources="showSourcesFor"
          />
        </div>

        <!-- Completed Section (if grouped and has checked items) -->
        <div v-if="checkedItems.length > 0 && groupByCategory" class="completed-section">
          <q-expansion-item
            label="Purchased"
            :caption="`${checkedItems.length} items`"
            icon="check_circle"
            header-class="completed-header"
          >
            <div class="completed-items">
              <ShoppingListItem
                v-for="item in checkedItems"
                :key="item.id"
                :item="item"
                :show-actions="true"
                @toggle="handleToggle"
                @delete="handleDeleteItem"
              />
            </div>
          </q-expansion-item>
        </div>
      </div>
    </div>

    <!-- FAB for adding items -->
    <q-page-sticky position="bottom-right" :offset="[24, 24]">
      <q-btn
        fab
        icon="add"
        color="primary"
        @click="showAddDialog = true"
      />
    </q-page-sticky>

    <!-- Add Item Dialog -->
    <AddItemDialog
      v-model="showAddDialog"
      :loading="addingItem"
      @add="handleAddItem"
    />

    <!-- Rename Dialog -->
    <q-dialog v-model="showRenameDialog">
      <q-card style="min-width: 350px; border-radius: 16px">
        <q-card-section>
          <div class="dialog-title">Rename List</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="newName"
            label="List Name"
            outlined
            dense
            autofocus
            @keyup.enter="handleRename"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-4 q-pb-4">
          <q-btn flat label="Cancel" color="grey" @click="showRenameDialog = false" />
          <q-btn
            unelevated
            label="Save"
            color="primary"
            :loading="renaming"
            :disable="!newName.trim()"
            @click="handleRename"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Sources Dialog -->
    <q-dialog v-model="showSourcesDialog">
      <q-card style="min-width: 320px; border-radius: 16px">
        <q-card-section>
          <div class="dialog-title">Recipes Using This Ingredient</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-list v-if="selectedItem?.sources?.length">
            <q-item v-for="source in selectedItem.sources" :key="source.recipeId">
              <q-item-section avatar>
                <q-icon name="restaurant" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ source.recipeName }}</q-item-label>
                <q-item-label caption v-if="source.quantity">
                  {{ source.quantity }} {{ source.unit || '' }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
          <p v-else class="tw-text-gray-500 tw-text-center">
            This is a custom item.
          </p>
        </q-card-section>

        <q-card-actions align="right" class="q-px-4 q-pb-4">
          <q-btn flat label="Close" color="grey" @click="showSourcesDialog = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useShoppinglistDetail } from '../composables/useShoppinglist';
import { useShoppinglistStore } from '../store/shoppinglistStore';
import ShoppingListItem from '../components/ShoppingListItem.vue';
import ShoppingListGroup from '../components/ShoppingListGroup.vue';
import AddItemDialog from '../components/AddItemDialog.vue';
import type { ShoppingListItemJSON } from '../types/shoppinglist';

const route = useRoute();
const router = useRouter();
const $q = useQuasar();
const store = useShoppinglistStore();

const listId = () => route.params.id as string;
const { list, loading, error, groupedItems, progress, checkedItems, toggleItem, addItem, deleteItem, updateList, refresh } = useShoppinglistDetail(listId);

// UI State
const groupByCategory = ref(true);
const showAddDialog = ref(false);
const showRenameDialog = ref(false);
const showSourcesDialog = ref(false);
const selectedItem = ref<ShoppingListItemJSON | null>(null);
const newName = ref('');
const addingItem = ref(false);
const renaming = ref(false);

// Initialize rename field when dialog opens
watch(showRenameDialog, (open) => {
  if (open && list.value) {
    newName.value = list.value.name;
  }
});

async function handleToggle(itemId: string) {
  try {
    await toggleItem(itemId);
  } catch {
    $q.notify({
      type: 'negative',
      message: 'Failed to update item',
    });
  }
}

async function handleAddItem(item: { customName: string; quantity?: number; unit?: string; notes?: string }) {
  addingItem.value = true;
  try {
    await addItem({
      customName: item.customName,
      quantity: item.quantity,
      unit: item.unit,
      notes: item.notes,
    });
    showAddDialog.value = false;
    $q.notify({
      type: 'positive',
      message: 'Item added',
    });
  } catch {
    $q.notify({
      type: 'negative',
      message: 'Failed to add item',
    });
  } finally {
    addingItem.value = false;
  }
}

function handleDeleteItem(itemId: string) {
  $q.dialog({
    title: 'Delete Item',
    message: 'Are you sure you want to remove this item?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    deleteItem(itemId)
      .then(() => {
        $q.notify({
          type: 'positive',
          message: 'Item removed',
        });
      })
      .catch(() => {
        $q.notify({
          type: 'negative',
          message: 'Failed to remove item',
        });
      });
  });
}

async function handleRename() {
  if (!newName.value.trim()) return;
  renaming.value = true;
  try {
    await updateList(newName.value.trim());
    showRenameDialog.value = false;
  } catch {
    $q.notify({
      type: 'negative',
      message: 'Failed to rename list',
    });
  } finally {
    renaming.value = false;
  }
}

function confirmDelete() {
  $q.dialog({
    title: 'Delete Shopping List',
    message: 'Are you sure you want to delete this shopping list? This cannot be undone.',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    store
      .deleteList(listId())
      .then(() => {
        $q.notify({
          type: 'positive',
          message: 'Shopping list deleted',
        });
        void router.push({ name: 'shopping-lists' });
      })
      .catch(() => {
        $q.notify({
          type: 'negative',
          message: 'Failed to delete list',
        });
      });
  });
}

function showSourcesFor(item: ShoppingListItemJSON) {
  selectedItem.value = item;
  showSourcesDialog.value = true;
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600;9..144,700&display=swap');

.shopping-list-detail-page {
  padding-top: env(safe-area-inset-top);
  background: linear-gradient(180deg, #fff8f5 0%, #ffffff 100%);
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%);
  padding: 12px 8px 16px;
  color: white;
}

.back-btn,
.menu-btn {
  color: white;
}

.page-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.progress-info {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  opacity: 0.9;
  margin-top: 2px;
}

.progress-bar-container {
  margin-top: 12px;
  padding: 0 8px;
}

.progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: white;
  border-radius: 2px;
  transition: width 0.3s ease;
}

.view-toggle {
  display: flex;
  justify-content: center;
}

.flat-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.completed-section {
  margin-top: 24px;
}

.completed-header {
  color: #9b918c;
}

.completed-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px 0;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
}

.dialog-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
}
</style>
