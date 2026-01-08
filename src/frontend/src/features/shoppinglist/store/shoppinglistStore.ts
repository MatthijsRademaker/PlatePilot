import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type {
  ShoppingListJSON,
  ShoppingListSummaryJSON,
  ShoppingListItemJSON,
  CreateFromRecipesRequest,
  AddItemRequest,
  UpdateItemRequest,
  GroupedItems,
  groupItemsByCategory,
} from '../types/shoppinglist';
import { groupItemsByCategory as groupItems } from '../types/shoppinglist';
import * as api from '../api/shoppinglistApi';

export const useShoppinglistStore = defineStore('shoppinglist', () => {
  // State
  const lists = ref<ShoppingListSummaryJSON[]>([]);
  const currentList = ref<ShoppingListJSON | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pageIndex = ref(1);
  const pageSize = ref(20);
  const totalCount = ref(0);

  // Getters
  const hasMore = computed(() => lists.value.length < totalCount.value);
  const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

  const uncheckedItems = computed(() =>
    currentList.value?.items.filter((item) => !item.checked) ?? []
  );

  const checkedItems = computed(() =>
    currentList.value?.items.filter((item) => item.checked) ?? []
  );

  const groupedItems = computed<GroupedItems[]>(() => {
    if (!currentList.value) return [];
    return groupItems(currentList.value.items);
  });

  const progress = computed(() => {
    if (!currentList.value || currentList.value.totalItems === 0) return 0;
    return (currentList.value.checkedItems / currentList.value.totalItems) * 100;
  });

  // Actions
  async function fetchLists(page = 1) {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.getShoppingLists({
        pageIndex: page,
        pageSize: pageSize.value,
      });
      lists.value = response.items ?? [];
      pageIndex.value = response.pageIndex ?? page;
      totalCount.value = response.totalCount ?? 0;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch shopping lists';
    } finally {
      loading.value = false;
    }
  }

  async function fetchListById(id: string) {
    loading.value = true;
    error.value = null;
    try {
      currentList.value = await api.getShoppingListById(id);
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch shopping list';
    } finally {
      loading.value = false;
    }
  }

  async function createFromRecipes(request: CreateFromRecipesRequest): Promise<ShoppingListJSON | null> {
    loading.value = true;
    error.value = null;
    try {
      const newList = await api.createFromRecipes(request);
      // Add to the beginning of the lists
      lists.value.unshift({
        id: newList.id,
        name: newList.name,
        totalItems: newList.totalItems,
        checkedItems: newList.checkedItems,
        createdAt: newList.createdAt,
        updatedAt: newList.updatedAt,
      });
      return newList;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create shopping list';
      return null;
    } finally {
      loading.value = false;
    }
  }

  async function createEmptyList(name?: string): Promise<ShoppingListJSON | null> {
    loading.value = true;
    error.value = null;
    try {
      const newList = await api.createShoppingList({ name });
      lists.value.unshift({
        id: newList.id,
        name: newList.name,
        totalItems: newList.totalItems,
        checkedItems: newList.checkedItems,
        createdAt: newList.createdAt,
        updatedAt: newList.updatedAt,
      });
      return newList;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create shopping list';
      return null;
    } finally {
      loading.value = false;
    }
  }

  async function updateList(id: string, name: string) {
    error.value = null;
    try {
      const updated = await api.updateShoppingList(id, { name });
      if (currentList.value?.id === id) {
        currentList.value = updated;
      }
      // Update in lists array
      const idx = lists.value.findIndex((l) => l.id === id);
      if (idx !== -1) {
        lists.value[idx].name = name;
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update shopping list';
      throw e;
    }
  }

  async function deleteList(id: string) {
    error.value = null;
    try {
      await api.deleteShoppingList(id);
      lists.value = lists.value.filter((l) => l.id !== id);
      if (currentList.value?.id === id) {
        currentList.value = null;
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete shopping list';
      throw e;
    }
  }

  async function addItem(listId: string, request: AddItemRequest) {
    error.value = null;
    try {
      const newItem = await api.addItem(listId, request);
      if (currentList.value?.id === listId) {
        currentList.value.items.push(newItem);
        currentList.value.totalItems++;
      }
      return newItem;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to add item';
      throw e;
    }
  }

  async function updateItem(listId: string, itemId: string, request: UpdateItemRequest) {
    error.value = null;
    try {
      await api.updateItem(listId, itemId, request);
      if (currentList.value?.id === listId) {
        const idx = currentList.value.items.findIndex((i) => i.id === itemId);
        if (idx !== -1) {
          const item = currentList.value.items[idx];
          if (request.quantity !== undefined) item.quantity = request.quantity;
          if (request.unit !== undefined) item.unit = request.unit;
          if (request.notes !== undefined) item.notes = request.notes;
          if (request.checked !== undefined) {
            const wasChecked = item.checked;
            item.checked = request.checked;
            if (wasChecked !== request.checked) {
              currentList.value.checkedItems += request.checked ? 1 : -1;
            }
          }
        }
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update item';
      throw e;
    }
  }

  async function toggleItem(listId: string, itemId: string) {
    error.value = null;
    // Optimistic update
    let previousState = false;
    if (currentList.value?.id === listId) {
      const idx = currentList.value.items.findIndex((i) => i.id === itemId);
      if (idx !== -1) {
        previousState = currentList.value.items[idx].checked;
        currentList.value.items[idx].checked = !previousState;
        currentList.value.checkedItems += previousState ? -1 : 1;
      }
    }

    try {
      const result = await api.toggleItemChecked(listId, itemId);
      // Sync with server state in case of race condition
      if (currentList.value?.id === listId) {
        const idx = currentList.value.items.findIndex((i) => i.id === itemId);
        if (idx !== -1) {
          const item = currentList.value.items[idx];
          if (item.checked !== result.checked) {
            currentList.value.checkedItems += result.checked ? 1 : -1;
            item.checked = result.checked;
          }
        }
      }
      return result.checked;
    } catch (e) {
      // Revert optimistic update
      if (currentList.value?.id === listId) {
        const idx = currentList.value.items.findIndex((i) => i.id === itemId);
        if (idx !== -1) {
          currentList.value.items[idx].checked = previousState;
          currentList.value.checkedItems += previousState ? 1 : -1;
        }
      }
      error.value = e instanceof Error ? e.message : 'Failed to toggle item';
      throw e;
    }
  }

  async function deleteItem(listId: string, itemId: string) {
    error.value = null;
    try {
      await api.deleteItem(listId, itemId);
      if (currentList.value?.id === listId) {
        const idx = currentList.value.items.findIndex((i) => i.id === itemId);
        if (idx !== -1) {
          if (currentList.value.items[idx].checked) {
            currentList.value.checkedItems--;
          }
          currentList.value.items.splice(idx, 1);
          currentList.value.totalItems--;
        }
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete item';
      throw e;
    }
  }

  function clearCurrentList() {
    currentList.value = null;
  }

  function clearError() {
    error.value = null;
  }

  return {
    // State
    lists,
    currentList,
    loading,
    error,
    pageIndex,
    pageSize,
    totalCount,
    // Getters
    hasMore,
    totalPages,
    uncheckedItems,
    checkedItems,
    groupedItems,
    progress,
    // Actions
    fetchLists,
    fetchListById,
    createFromRecipes,
    createEmptyList,
    updateList,
    deleteList,
    addItem,
    updateItem,
    toggleItem,
    deleteItem,
    clearCurrentList,
    clearError,
  };
});
