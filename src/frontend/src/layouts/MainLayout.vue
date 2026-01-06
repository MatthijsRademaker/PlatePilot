<template>
  <q-layout view="hHh Lpr fFf">
    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="main-footer">
      <q-tabs align="justify" class="footer-tabs" indicator-color="transparent">
        <q-route-tab
          v-for="tab in navLinks"
          :key="tab.to"
          :to="{ name: tab.to }"
          :icon="tab.icon"
          :label="tab.label"
          class="footer-tab"
        />
      </q-tabs>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
interface NavLink {
  to: string;
  icon: string;
  label: string;
}

const navLinks: NavLink[] = [
  {
    to: 'home',
    icon: 'mdi-home',
    label: 'Home',
  },
  {
    to: 'recipes',
    icon: 'mdi-book-open-variant-outline',
    label: 'Recipes',
  },
  {
    to: 'mealplan',
    icon: 'mdi-calendar',
    label: 'Meal Plan',
  },
  {
    to: 'search',
    icon: 'mdi-magnify',
    label: 'Search',
  },
];
</script>

<style scoped lang="scss">
.main-footer {
  background: #2d1f1a;
}

.footer-tabs {
  :deep(.q-tab) {
    color: rgba(255, 255, 255, 0.7);
  }

  :deep(.q-tab--active) {
    color: #ff7f50;
  }
}

.footer-tab {
  min-height: 56px;
}
</style>
