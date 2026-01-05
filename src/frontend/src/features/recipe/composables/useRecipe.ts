import { ref, computed } from 'vue';
import {
  useGetRecipeAll,
  useGetRecipeId,
  useGetRecipeSimilar,
} from '@/api/generated/platepilot';

export function useRecipeList() {
  const pageIndex = ref(1);
  const pageSize = ref(20);

  const { data, isPending, error, refetch } = useGetRecipeAll(
    computed(() => ({
      pageIndex: pageIndex.value,
      pageSize: pageSize.value,
    })),
  );

  const recipes = computed(() => data.value?.items ?? []);
  const totalCount = computed(() => data.value?.totalCount ?? 0);
  const totalPages = computed(() => data.value?.totalPages ?? 0);
  const hasMore = computed(() => recipes.value.length < totalCount.value);

  function loadPage(page: number) {
    pageIndex.value = page;
  }

  function refresh() {
    pageIndex.value = 1;
    void refetch();
  }

  return {
    recipes,
    loading: isPending,
    error: computed(() => error.value?.error ?? null),
    pageIndex,
    totalPages,
    hasMore,
    loadPage,
    refresh,
  };
}

export function useRecipeDetail(recipeId: () => string) {
  const { data, isPending, error } = useGetRecipeId(computed(() => recipeId()));

  return {
    recipe: data,
    loading: isPending,
    error: computed(() => error.value?.message ?? null),
  };
}

export function useSimilarRecipes(recipeId: () => string, amount = 5) {
  const { data, isPending, error } = useGetRecipeSimilar(
    computed(() => ({
      recipe: recipeId(),
      amount,
    })),
  );

  return {
    recipes: computed(() => data.value ?? []),
    loading: isPending,
    error: computed(() => error.value?.error ?? null),
  };
}
