<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="min-width: 350px; border-radius: 16px">
      <q-card-section>
        <div class="dialog-title">Add Item</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input
          v-model="itemName"
          label="Item name"
          autofocus
          outlined
          dense
          class="tw-mb-4"
          @keyup.enter="handleAdd"
        />

        <div class="tw-flex tw-gap-3">
          <q-input
            v-model.number="quantity"
            label="Quantity"
            type="number"
            outlined
            dense
            class="tw-flex-1"
            min="0"
            step="0.5"
          />

          <q-select
            v-model="unit"
            label="Unit"
            :options="unitOptions"
            outlined
            dense
            class="tw-flex-1"
            clearable
          />
        </div>

        <q-input
          v-model="notes"
          label="Notes (optional)"
          outlined
          dense
          class="tw-mt-4"
        />
      </q-card-section>

      <q-card-actions align="right" class="q-px-4 q-pb-4">
        <q-btn flat label="Cancel" color="grey" @click="handleCancel" />
        <q-btn
          unelevated
          label="Add"
          color="primary"
          :loading="loading"
          :disable="!itemName.trim()"
          @click="handleAdd"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  add: [item: { customName: string; quantity?: number; unit?: string; notes?: string }];
}>();

const itemName = ref('');
const quantity = ref<number | undefined>(undefined);
const unit = ref<string | undefined>(undefined);
const notes = ref('');

const unitOptions = [
  'pcs',
  'g',
  'kg',
  'ml',
  'l',
  'oz',
  'lb',
  'cup',
  'tbsp',
  'tsp',
  'bunch',
  'can',
  'pack',
];

// Reset form when dialog opens
watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      itemName.value = '';
      quantity.value = undefined;
      unit.value = undefined;
      notes.value = '';
    }
  }
);

function handleAdd() {
  if (!itemName.value.trim()) return;

  emit('add', {
    customName: itemName.value.trim(),
    quantity: quantity.value,
    unit: unit.value,
    notes: notes.value.trim() || undefined,
  });
}

function handleCancel() {
  emit('update:modelValue', false);
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,600&display=swap');

.dialog-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  font-weight: 600;
  color: #2d1f1a;
}
</style>
