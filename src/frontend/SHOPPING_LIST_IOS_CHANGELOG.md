# Shopping List Feature - iOS Implementation Guide

**Feature Branch:** `feature/shopping-list`
**Priority:** High
**Status:** Backend + Web Frontend Complete

## Overview

The Shopping List feature allows users to generate grocery lists from their meal plan recipes. Ingredients are automatically aggregated (e.g., "2 eggs" from Recipe A + "3 eggs" from Recipe B = "5 eggs"), grouped by category for easy shopping, and can be checked off as purchased.

## API Endpoints

All endpoints require authentication. Base URL: `/v1/shoppinglist`

### 1. List Shopping Lists (Paginated)

```
GET /v1/shoppinglist/
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `pageIndex` | int | 1 | Page number (1-indexed) |
| `pageSize` | int | 20 | Items per page |

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "Meal Plan - Week of Jan 6",
      "totalItems": 15,
      "checkedItems": 3,
      "createdAt": "2025-01-06T10:00:00Z",
      "updatedAt": "2025-01-06T12:30:00Z",
      "completedAt": null
    }
  ],
  "pageIndex": 1,
  "pageSize": 20,
  "totalCount": 5,
  "totalPages": 1
}
```

---

### 2. Get Shopping List Detail

```
GET /v1/shoppinglist/{id}
```

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "name": "Meal Plan - Week of Jan 6",
  "weekStartDate": "2025-01-06",
  "items": [
    {
      "id": "item-uuid",
      "ingredient": {
        "id": "ingredient-uuid",
        "name": "Eggs"
      },
      "category": {
        "id": "category-uuid",
        "name": "Dairy & Eggs",
        "displayOrder": 2
      },
      "customName": null,
      "quantity": 5,
      "unit": "large",
      "displayQuantity": "5 large",
      "checked": false,
      "notes": null,
      "isCustom": false,
      "sources": [
        {
          "recipeId": "recipe-uuid-1",
          "recipeName": "Scrambled Eggs",
          "quantity": 2,
          "unit": "large"
        },
        {
          "recipeId": "recipe-uuid-2",
          "recipeName": "French Toast",
          "quantity": 3,
          "unit": "large"
        }
      ],
      "createdAt": "2025-01-06T10:00:00Z"
    }
  ],
  "recipes": [
    { "id": "recipe-uuid-1", "name": "Scrambled Eggs" },
    { "id": "recipe-uuid-2", "name": "French Toast" }
  ],
  "totalItems": 15,
  "checkedItems": 3,
  "createdAt": "2025-01-06T10:00:00Z",
  "updatedAt": "2025-01-06T12:30:00Z",
  "completedAt": null
}
```

**Error:** `404 Not Found` if list doesn't exist or belongs to another user

---

### 3. Create Shopping List from Recipes

```
POST /v1/shoppinglist/from-recipes
```

**Request Body:**
```json
{
  "name": "Meal Plan - Week of Jan 6",
  "recipeIds": ["uuid-1", "uuid-2", "uuid-3"],
  "weekStartDate": "2025-01-06"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | No | Auto-generates if not provided |
| `recipeIds` | string[] | Yes | Array of recipe UUIDs |
| `weekStartDate` | string | No | ISO date for meal plan week |

**Response:** `201 Created` - Returns full `ShoppingListJSON` (same as detail endpoint)

**Notes:**
- Ingredients from all recipes are aggregated
- Same ingredient with same unit = quantities summed
- Same ingredient with different units = kept separate
- Each item tracks which recipes it came from (`sources` array)

---

### 4. Create Empty Shopping List

```
POST /v1/shoppinglist/
```

**Request Body:**
```json
{
  "name": "Quick Trip"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | No | Defaults to "Shopping List" |

**Response:** `201 Created` - Returns full `ShoppingListJSON`

---

### 5. Update Shopping List

```
PATCH /v1/shoppinglist/{id}
```

**Request Body:**
```json
{
  "name": "Updated Name"
}
```

**Response:** `200 OK` - Returns updated `ShoppingListJSON`

---

### 6. Delete Shopping List

```
DELETE /v1/shoppinglist/{id}
```

**Response:** `204 No Content`

---

### 7. Add Item to List

```
POST /v1/shoppinglist/{listId}/items
```

**Request Body:**
```json
{
  "ingredientId": "ingredient-uuid",
  "customName": null,
  "quantity": 2,
  "unit": "cups",
  "notes": "Get organic if available"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ingredientId` | string | No* | Link to existing ingredient |
| `customName` | string | No* | Free-text item name |
| `quantity` | number | No | Amount needed |
| `unit` | string | No | Unit of measurement |
| `notes` | string | No | Additional notes |

*Either `ingredientId` OR `customName` should be provided. If `customName` is set, `isCustom` will be `true`.

**Response:** `201 Created`
```json
{
  "id": "new-item-uuid",
  "ingredient": null,
  "category": null,
  "customName": "Fancy Cheese",
  "quantity": 1,
  "unit": "block",
  "displayQuantity": "1 block",
  "checked": false,
  "notes": "Get organic if available",
  "isCustom": true,
  "sources": [],
  "createdAt": "2025-01-06T14:00:00Z"
}
```

---

### 8. Update Item

```
PATCH /v1/shoppinglist/{listId}/items/{itemId}
```

**Request Body:**
```json
{
  "quantity": 3,
  "unit": "cups",
  "notes": "Updated notes",
  "checked": true
}
```

All fields are optional. Only provided fields are updated.

**Response:** `200 OK`
```json
{
  "id": "item-uuid"
}
```

---

### 9. Toggle Item Checked State

```
POST /v1/shoppinglist/{listId}/items/{itemId}/toggle
```

**Request Body:** None

**Response:** `200 OK`
```json
{
  "checked": true
}
```

**Notes:**
- Use optimistic UI updates for better UX
- Toggle the local state immediately, then sync with server
- Revert on error

---

### 10. Delete Item

```
DELETE /v1/shoppinglist/{listId}/items/{itemId}
```

**Response:** `204 No Content`

---

## Swift Data Models

```swift
import Foundation

// MARK: - Shopping List Models

struct ShoppingList: Codable, Identifiable {
    let id: String
    var name: String
    let weekStartDate: String?
    var items: [ShoppingListItem]
    let recipes: [RecipeRef]
    var totalItems: Int
    var checkedItems: Int
    let createdAt: Date
    var updatedAt: Date
    let completedAt: Date?

    var progress: Double {
        guard totalItems > 0 else { return 0 }
        return Double(checkedItems) / Double(totalItems)
    }

    var isCompleted: Bool {
        totalItems > 0 && checkedItems == totalItems
    }
}

struct ShoppingListSummary: Codable, Identifiable {
    let id: String
    let name: String
    let totalItems: Int
    let checkedItems: Int
    let createdAt: Date
    let updatedAt: Date
    let completedAt: Date?

    var progress: Double {
        guard totalItems > 0 else { return 0 }
        return Double(checkedItems) / Double(totalItems)
    }
}

struct ShoppingListItem: Codable, Identifiable {
    let id: String
    let ingredient: Ingredient?
    let category: IngredientCategory?
    let customName: String?
    var quantity: Double?
    var unit: String?
    let displayQuantity: String
    var checked: Bool
    var notes: String?
    let isCustom: Bool
    let sources: [ItemSource]?
    let createdAt: Date

    /// Display name (ingredient name or custom name)
    var displayName: String {
        customName ?? ingredient?.name ?? "Unknown Item"
    }
}

struct Ingredient: Codable, Identifiable {
    let id: String
    let name: String
}

struct IngredientCategory: Codable, Identifiable {
    let id: String
    let name: String
    let displayOrder: Int
}

struct ItemSource: Codable {
    let recipeId: String
    let recipeName: String
    let quantity: Double?
    let unit: String?
}

struct RecipeRef: Codable, Identifiable {
    let id: String
    let name: String
}

// MARK: - Paginated Response

struct PaginatedShoppingLists: Codable {
    let items: [ShoppingListSummary]
    let pageIndex: Int
    let pageSize: Int
    let totalCount: Int
    let totalPages: Int
}

// MARK: - Request Models

struct CreateFromRecipesRequest: Codable {
    var name: String?
    let recipeIds: [String]
    var weekStartDate: String?
}

struct CreateShoppingListRequest: Codable {
    var name: String?
}

struct UpdateShoppingListRequest: Codable {
    let name: String
}

struct AddItemRequest: Codable {
    var ingredientId: String?
    var customName: String?
    var quantity: Double?
    var unit: String?
    var notes: String?
}

struct UpdateItemRequest: Codable {
    var quantity: Double?
    var unit: String?
    var notes: String?
    var checked: Bool?
}

struct ToggleResponse: Codable {
    let checked: Bool
}

// MARK: - Grouped Items (for UI)

struct GroupedItems: Identifiable {
    let id: String // categoryId
    let categoryName: String
    let displayOrder: Int
    var items: [ShoppingListItem]
}

extension Array where Element == ShoppingListItem {
    /// Groups items by category for display
    func groupedByCategory() -> [GroupedItems] {
        var groups: [String: GroupedItems] = [:]
        var uncategorized: [ShoppingListItem] = []

        for item in self where !item.checked {
            if let category = item.category {
                if groups[category.id] == nil {
                    groups[category.id] = GroupedItems(
                        id: category.id,
                        categoryName: category.name,
                        displayOrder: category.displayOrder,
                        items: []
                    )
                }
                groups[category.id]?.items.append(item)
            } else {
                uncategorized.append(item)
            }
        }

        var result = groups.values.sorted { $0.displayOrder < $1.displayOrder }

        if !uncategorized.isEmpty {
            result.append(GroupedItems(
                id: "uncategorized",
                categoryName: "Other",
                displayOrder: 999,
                items: uncategorized
            ))
        }

        return result
    }
}
```

---

## UI/UX Patterns

### Navigation Structure

```
Tab Bar
├── Home
├── Recipes
├── Plan (Meal Plan)
├── Shop (Shopping Lists) ← NEW
└── Search
```

### Shopping Lists Page

- List of shopping list cards
- Each card shows: name, progress bar, "X of Y items"
- Pull-to-refresh
- FAB to create new empty list
- Swipe-to-delete on cards

### Shopping List Detail Page

**Header:**
- Back button
- List name (editable)
- Progress indicator ("5 of 12 items")
- Overflow menu (Rename, Delete)
- Progress bar

**View Toggle:**
- Grouped (by category) - default
- All Items (flat list)

**Item Row:**
- Checkbox (large tap target)
- Item name (ingredient name or custom name)
- Quantity + unit
- Info button (shows source recipes)
- Swipe-to-delete

**Completed Section:**
- Collapsible section for checked items
- Shows at bottom: "Purchased (X items)"

**FAB:**
- Add custom item

### Optimistic Updates

For checkbox toggling, implement optimistic updates:

```swift
func toggleItem(_ item: ShoppingListItem) async {
    // 1. Immediately update local state
    let previousState = item.checked
    updateLocalState(itemId: item.id, checked: !previousState)

    do {
        // 2. Call API
        let response = try await api.toggleItem(listId: listId, itemId: item.id)
        // 3. Sync with server state (in case of race condition)
        updateLocalState(itemId: item.id, checked: response.checked)
    } catch {
        // 4. Revert on error
        updateLocalState(itemId: item.id, checked: previousState)
        showError("Failed to update item")
    }
}
```

### Meal Plan Integration

On the Meal Plan page, add a "Generate Shopping List" button:

```swift
func generateShoppingList() async {
    // 1. Get unique recipe IDs from current week
    let recipeIds = currentWeek.days
        .flatMap { $0.slots }
        .compactMap { $0.recipe?.id }
        .unique()

    guard !recipeIds.isEmpty else {
        showWarning("No recipes in your meal plan. Add some recipes first!")
        return
    }

    // 2. Create shopping list
    do {
        let request = CreateFromRecipesRequest(
            name: "Meal Plan - \(currentWeek.label)",
            recipeIds: recipeIds
        )
        let list = try await api.createFromRecipes(request)

        showSuccess("Shopping list created with \(list.totalItems) items")
        navigate(to: .shoppingListDetail(id: list.id))
    } catch {
        showError("Failed to generate shopping list")
    }
}
```

---

## Ingredient Categories

The backend seeds these default categories (ordered for typical grocery store layout):

| Order | Category |
|-------|----------|
| 1 | Produce |
| 2 | Dairy & Eggs |
| 3 | Meat & Seafood |
| 4 | Bakery |
| 5 | Pantry |
| 6 | Frozen |
| 7 | Beverages |
| 8 | Condiments & Sauces |
| 9 | Spices & Seasonings |
| 10 | Other |

---

## Error Handling

| Status Code | Meaning | Action |
|-------------|---------|--------|
| 400 | Bad Request | Show validation error |
| 401 | Unauthorized | Redirect to login |
| 404 | Not Found | Show "List not found" |
| 500 | Server Error | Show generic error, allow retry |

---

## Testing Checklist

- [ ] Create shopping list from meal plan with multiple recipes
- [ ] Verify ingredient aggregation (same ingredient, same unit)
- [ ] Verify items are grouped by category
- [ ] Toggle item checked state (optimistic update)
- [ ] Add custom item
- [ ] Delete item
- [ ] Rename list
- [ ] Delete list
- [ ] View item sources (which recipes use this ingredient)
- [ ] Empty state when no items
- [ ] Loading states
- [ ] Error states with retry
- [ ] Pull-to-refresh on list page

---

## Files Changed (Web Frontend Reference)

```
src/features/shoppinglist/
├── types/shoppinglist.ts          # Data models
├── api/shoppinglistApi.ts         # API calls
├── store/shoppinglistStore.ts     # State management
├── composables/useShoppinglist.ts # Reusable logic
├── components/
│   ├── ShoppingListCard.vue       # List preview card
│   ├── ShoppingListItem.vue       # Checkbox item row
│   ├── ShoppingListGroup.vue      # Category section
│   └── AddItemDialog.vue          # Add item modal
├── pages/
│   ├── ShoppingListsPage.vue      # All lists
│   └── ShoppingListDetailPage.vue # Single list detail
└── routes.ts

src/features/mealplan/pages/MealplanPage.vue  # Added generate button
src/layouts/MainLayout.vue                      # Added nav tab
```

---

## Questions?

Reach out in #mobile-dev or check the web implementation for reference.
