/**
 * Shopping List Types
 * These types match the backend BFF handler JSON responses
 */

export interface IngredientJSON {
  id: string;
  name: string;
}

export interface IngredientCategoryJSON {
  id: string;
  name: string;
  displayOrder: number;
}

export interface ItemSourceJSON {
  recipeId: string;
  recipeName: string;
  quantity?: number;
  unit?: string;
}

export interface RecipeRefJSON {
  id: string;
  name: string;
}

export interface ShoppingListItemJSON {
  id: string;
  ingredient?: IngredientJSON;
  category?: IngredientCategoryJSON;
  customName?: string;
  quantity?: number;
  unit?: string;
  displayQuantity: string;
  checked: boolean;
  notes?: string;
  isCustom: boolean;
  sources?: ItemSourceJSON[];
  createdAt: string;
}

export interface ShoppingListJSON {
  id: string;
  name: string;
  weekStartDate?: string;
  items: ShoppingListItemJSON[];
  recipes: RecipeRefJSON[];
  totalItems: number;
  checkedItems: number;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
}

export interface ShoppingListSummaryJSON {
  id: string;
  name: string;
  totalItems: number;
  checkedItems: number;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
}

export interface PaginatedShoppingListsJSON {
  items: ShoppingListSummaryJSON[];
  pageIndex: number;
  pageSize: number;
  totalCount: number;
  totalPages: number;
}

// Request types
export interface CreateFromRecipesRequest {
  name?: string;
  recipeIds: string[];
  weekStartDate?: string;
}

export interface CreateShoppingListRequest {
  name?: string;
}

export interface UpdateShoppingListRequest {
  name: string;
}

export interface AddItemRequest {
  ingredientId?: string;
  customName?: string;
  quantity?: number;
  unit?: string;
  notes?: string;
}

export interface UpdateItemRequest {
  quantity?: number;
  unit?: string;
  notes?: string;
  checked?: boolean;
}

// Grouped items by category (for UI)
export interface GroupedItems {
  categoryId: string;
  categoryName: string;
  displayOrder: number;
  items: ShoppingListItemJSON[];
}

// Helper function to group items by category
export function groupItemsByCategory(items: ShoppingListItemJSON[]): GroupedItems[] {
  const groups: Map<string, GroupedItems> = new Map();
  const uncategorized: ShoppingListItemJSON[] = [];

  for (const item of items) {
    if (item.category) {
      const key = item.category.id;
      if (!groups.has(key)) {
        groups.set(key, {
          categoryId: item.category.id,
          categoryName: item.category.name,
          displayOrder: item.category.displayOrder,
          items: [],
        });
      }
      groups.get(key)!.items.push(item);
    } else {
      uncategorized.push(item);
    }
  }

  // Convert to array and sort by display order
  const result = Array.from(groups.values()).sort(
    (a, b) => a.displayOrder - b.displayOrder
  );

  // Add uncategorized at the end if any
  if (uncategorized.length > 0) {
    result.push({
      categoryId: 'uncategorized',
      categoryName: 'Other',
      displayOrder: 999,
      items: uncategorized,
    });
  }

  return result;
}
