<template>
  <div class="shopping-list-group">
    <div v-if="showHeader" class="group-header">
      <span class="group-name">{{ group.categoryName }}</span>
      <span class="group-count">{{ group.items.length }}</span>
    </div>

    <div class="group-items">
      <ShoppingListItem
        v-for="item in group.items"
        :key="item.id"
        :item="item"
        :show-actions="showItemActions"
        @toggle="$emit('toggle-item', $event)"
        @edit="$emit('edit-item', $event)"
        @delete="$emit('delete-item', $event)"
        @view-sources="$emit('view-sources', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GroupedItems, ShoppingListItemJSON } from '../types/shoppinglist';
import ShoppingListItem from './ShoppingListItem.vue';

defineProps<{
  group: GroupedItems;
  showHeader?: boolean;
  showItemActions?: boolean;
}>();

defineEmits<{
  'toggle-item': [itemId: string];
  'edit-item': [item: ShoppingListItemJSON];
  'delete-item': [itemId: string];
  'view-sources': [item: ShoppingListItemJSON];
}>();
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,500;9..40,600&display=swap');

.shopping-list-group {
  margin-bottom: 24px;
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px 8px;
}

.group-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 600;
  color: #6b5f5a;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.group-count {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 500;
  color: #9b918c;
  background: #f0ece9;
  padding: 2px 8px;
  border-radius: 10px;
}

.group-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
</style>
