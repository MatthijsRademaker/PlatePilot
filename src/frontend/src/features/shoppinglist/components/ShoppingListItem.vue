<template>
  <div
    class="shopping-list-item"
    :class="{ checked: item.checked }"
    @click="$emit('toggle', item.id)"
  >
    <q-checkbox
      :model-value="item.checked"
      dense
      class="item-checkbox"
      @update:model-value="$emit('toggle', item.id)"
      @click.stop
    />

    <div class="item-content">
      <span class="item-name">
        {{ item.ingredient?.name || item.customName || 'Unknown item' }}
      </span>
      <span class="item-quantity">{{ item.displayQuantity }}</span>
    </div>

    <q-btn
      v-if="showActions"
      icon="more_vert"
      flat
      round
      dense
      size="sm"
      class="action-btn"
      @click.stop="showMenu = !showMenu"
    >
      <q-menu v-model="showMenu">
        <q-list dense>
          <q-item clickable v-close-popup @click="$emit('edit', item)">
            <q-item-section avatar>
              <q-icon name="edit" size="18px" />
            </q-item-section>
            <q-item-section>Edit</q-item-section>
          </q-item>
          <q-item clickable v-close-popup @click="$emit('view-sources', item)">
            <q-item-section avatar>
              <q-icon name="info" size="18px" />
            </q-item-section>
            <q-item-section>View recipes</q-item-section>
          </q-item>
          <q-separator />
          <q-item clickable v-close-popup @click="$emit('delete', item.id)">
            <q-item-section avatar>
              <q-icon name="delete" size="18px" color="negative" />
            </q-item-section>
            <q-item-section class="text-negative">Delete</q-item-section>
          </q-item>
        </q-list>
      </q-menu>
    </q-btn>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { ShoppingListItemJSON } from '../types/shoppinglist';

defineProps<{
  item: ShoppingListItemJSON;
  showActions?: boolean;
}>();

defineEmits<{
  toggle: [itemId: string];
  edit: [item: ShoppingListItemJSON];
  delete: [itemId: string];
  'view-sources': [item: ShoppingListItemJSON];
}>();

const showMenu = ref(false);
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500&display=swap');

.shopping-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    opacity 0.2s ease;

  &:active {
    background: #f9f7f6;
  }

  &.checked {
    opacity: 0.6;

    .item-name {
      text-decoration: line-through;
      color: #9b918c;
    }
  }
}

.item-checkbox {
  flex-shrink: 0;
}

.item-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-width: 0;
  gap: 12px;
}

.item-name {
  font-family: 'DM Sans', sans-serif;
  font-size: 15px;
  font-weight: 500;
  color: #2d1f1a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-quantity {
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  font-weight: 400;
  color: #6b5f5a;
  white-space: nowrap;
  flex-shrink: 0;
}

.action-btn {
  flex-shrink: 0;
  opacity: 0.5;
  transition: opacity 0.15s ease;

  &:hover {
    opacity: 1;
  }
}

:deep(.q-checkbox__inner--falsy) {
  color: #c9c3bf;
}

:deep(.q-checkbox__inner--truthy) {
  color: #ff6347;
}
</style>
