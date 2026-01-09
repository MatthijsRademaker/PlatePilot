<template>
  <div class="shopping-list-card" @click="$emit('click', shoppingList)">
    <div class="card-header">
      <h3 class="card-title">{{ shoppingList.name }}</h3>
      <div class="item-count">
        <q-icon name="check_circle" size="14px" />
        {{ shoppingList.checkedItems }}/{{ shoppingList.totalItems }}
      </div>
    </div>

    <div class="progress-bar">
      <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
    </div>

    <div class="card-footer">
      <span class="date">{{ formattedDate }}</span>
      <span v-if="isCompleted" class="completed-badge">
        <q-icon name="done_all" size="12px" />
        Done
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { ShoppingListSummaryJSON } from '../types/shoppinglist';

const props = defineProps<{
  shoppingList: ShoppingListSummaryJSON;
}>();

defineEmits<{
  click: [shoppingList: ShoppingListSummaryJSON];
}>();

const progressPercent = computed(() => {
  if (props.shoppingList.totalItems === 0) return 0;
  return (props.shoppingList.checkedItems / props.shoppingList.totalItems) * 100;
});

const isCompleted = computed(() => {
  return (
    props.shoppingList.totalItems > 0 &&
    props.shoppingList.checkedItems === props.shoppingList.totalItems
  );
});

const formattedDate = computed(() => {
  const date = new Date(props.shoppingList.createdAt);
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
});
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Fraunces:opsz,wght@9..144,600&display=swap');

.shopping-list-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
  border: 1px solid rgba(45, 31, 26, 0.04);
  box-shadow: 0 2px 12px rgba(45, 31, 26, 0.04);

  &:active {
    transform: scale(0.98);
  }

  @media (hover: hover) {
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(255, 127, 80, 0.12);
    }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.card-title {
  font-family: 'Fraunces', serif;
  font-size: 16px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0;
  flex: 1;
  line-height: 1.3;
  padding-right: 12px;
}

.item-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #ff6347;
  background: #fff5f2;
  padding: 4px 10px;
  border-radius: 8px;
  white-space: nowrap;
}

.progress-bar {
  height: 4px;
  background: #f0ece9;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 12px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #ff7f50 0%, #ff6347 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.date {
  font-family: 'DM Sans', sans-serif;
  font-size: 12px;
  color: #9b918c;
}

.completed-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  font-weight: 600;
  color: #22c55e;
  background: #f0fdf4;
  padding: 3px 8px;
  border-radius: 6px;
}
</style>
