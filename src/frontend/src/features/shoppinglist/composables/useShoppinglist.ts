import { ref, computed, onMounted, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useShoppinglistStore } from '../store/shoppinglistStore';
import type { CreateFromRecipesRequest, AddItemRequest } from '../types/shoppinglist';

/**
 * Composable for managing the shopping lists page
 */
export function useShoppinglistList() {
  const store = useShoppinglistStore();
  const { lists, loading, error, pageIndex, totalPages, hasMore } = storeToRefs(store);

  onMounted(() => {
    if (lists.value.length === 0) {
      store.fetchLists();
    }
  });

  function loadPage(page: number) {
    store.fetchLists(page);
  }

  function refresh() {
    store.fetchLists(1);
  }

  return {
    lists,
    loading,
    error,
    pageIndex,
    totalPages,
    hasMore,
    loadPage,
    refresh,
    createFromRecipes: store.createFromRecipes,
    createEmptyList: store.createEmptyList,
    deleteList: store.deleteList,
  };
}

/**
 * Composable for managing a single shopping list detail page
 */
export function useShoppinglistDetail(listId: () => string) {
  const store = useShoppinglistStore();
  const { currentList, loading, error, groupedItems, progress, uncheckedItems, checkedItems } =
    storeToRefs(store);

  // Track the ID we've loaded
  const loadedId = ref<string | null>(null);

  // Load on mount and when ID changes
  function loadList() {
    const id = listId();
    if (id && id !== loadedId.value) {
      loadedId.value = id;
      store.fetchListById(id);
    }
  }

  onMounted(loadList);
  watch(listId, loadList);

  return {
    list: currentList,
    loading,
    error,
    groupedItems,
    progress,
    uncheckedItems,
    checkedItems,
    toggleItem: (itemId: string) => store.toggleItem(listId(), itemId),
    addItem: (request: AddItemRequest) => store.addItem(listId(), request),
    deleteItem: (itemId: string) => store.deleteItem(listId(), itemId),
    updateList: (name: string) => store.updateList(listId(), name),
    refresh: () => {
      loadedId.value = null;
      loadList();
    },
  };
}

/**
 * Composable for creating a shopping list from recipes
 * Used in the meal plan integration
 */
export function useCreateShoppingList() {
  const store = useShoppinglistStore();
  const creating = ref(false);
  const createError = ref<string | null>(null);

  async function createFromRecipes(request: CreateFromRecipesRequest) {
    creating.value = true;
    createError.value = null;
    try {
      const list = await store.createFromRecipes(request);
      return list;
    } catch (e) {
      createError.value = e instanceof Error ? e.message : 'Failed to create shopping list';
      throw e;
    } finally {
      creating.value = false;
    }
  }

  async function createEmpty(name?: string) {
    creating.value = true;
    createError.value = null;
    try {
      const list = await store.createEmptyList(name);
      return list;
    } catch (e) {
      createError.value = e instanceof Error ? e.message : 'Failed to create shopping list';
      throw e;
    } finally {
      creating.value = false;
    }
  }

  return {
    creating,
    error: createError,
    createFromRecipes,
    createEmpty,
  };
}
