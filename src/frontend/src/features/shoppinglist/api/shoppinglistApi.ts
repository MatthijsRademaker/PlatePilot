/**
 * Shopping List API
 * Manually defined API calls for shopping list feature
 * (These will be auto-generated once the backend swagger spec is updated)
 */

import { customInstance } from '@/api/mutator/custom-instance';
import type {
  ShoppingListJSON,
  PaginatedShoppingListsJSON,
  CreateFromRecipesRequest,
  CreateShoppingListRequest,
  UpdateShoppingListRequest,
  AddItemRequest,
  UpdateItemRequest,
  ShoppingListItemJSON,
} from '../types/shoppinglist';

/**
 * Get all shopping lists (paginated)
 */
export async function getShoppingLists(params?: {
  pageIndex?: number;
  pageSize?: number;
}): Promise<PaginatedShoppingListsJSON> {
  return customInstance<PaginatedShoppingListsJSON>({
    url: '/shoppinglist/',
    method: 'GET',
    params: params,
  });
}

/**
 * Get a shopping list by ID
 */
export async function getShoppingListById(id: string): Promise<ShoppingListJSON> {
  return customInstance<ShoppingListJSON>({
    url: `/shoppinglist/${id}`,
    method: 'GET',
  });
}

/**
 * Create a shopping list from recipe IDs
 */
export async function createFromRecipes(
  request: CreateFromRecipesRequest
): Promise<ShoppingListJSON> {
  return customInstance<ShoppingListJSON>({
    url: '/shoppinglist/from-recipes',
    method: 'POST',
    data: request,
  });
}

/**
 * Create an empty shopping list
 */
export async function createShoppingList(
  request: CreateShoppingListRequest
): Promise<ShoppingListJSON> {
  return customInstance<ShoppingListJSON>({
    url: '/shoppinglist/',
    method: 'POST',
    data: request,
  });
}

/**
 * Update a shopping list
 */
export async function updateShoppingList(
  id: string,
  request: UpdateShoppingListRequest
): Promise<ShoppingListJSON> {
  return customInstance<ShoppingListJSON>({
    url: `/shoppinglist/${id}`,
    method: 'PATCH',
    data: request,
  });
}

/**
 * Delete a shopping list
 */
export async function deleteShoppingList(id: string): Promise<void> {
  return customInstance<void>({
    url: `/shoppinglist/${id}`,
    method: 'DELETE',
  });
}

/**
 * Add an item to a shopping list
 */
export async function addItem(
  listId: string,
  request: AddItemRequest
): Promise<ShoppingListItemJSON> {
  return customInstance<ShoppingListItemJSON>({
    url: `/shoppinglist/${listId}/items`,
    method: 'POST',
    data: request,
  });
}

/**
 * Update an item in a shopping list
 */
export async function updateItem(
  listId: string,
  itemId: string,
  request: UpdateItemRequest
): Promise<{ id: string }> {
  return customInstance<{ id: string }>({
    url: `/shoppinglist/${listId}/items/${itemId}`,
    method: 'PATCH',
    data: request,
  });
}

/**
 * Toggle item checked state
 */
export async function toggleItemChecked(
  listId: string,
  itemId: string
): Promise<{ checked: boolean }> {
  return customInstance<{ checked: boolean }>({
    url: `/shoppinglist/${listId}/items/${itemId}/toggle`,
    method: 'POST',
  });
}

/**
 * Delete an item from a shopping list
 */
export async function deleteItem(listId: string, itemId: string): Promise<void> {
  return customInstance<void>({
    url: `/shoppinglist/${listId}/items/${itemId}`,
    method: 'DELETE',
  });
}
