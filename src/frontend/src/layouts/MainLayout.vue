<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-primary">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <router-link to="/" class="text-white text-decoration-none">
            PlatePilot
          </router-link>
        </q-toolbar-title>

        <q-btn flat round icon="search" @click="$router.push({ name: 'search' })" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header>Navigation</q-item-label>

        <q-item
          v-for="link in navLinks"
          :key="link.to"
          :to="{ name: link.to }"
          clickable
          v-ripple
          :active="$route.name === link.to"
          active-class="bg-primary text-white"
        >
          <q-item-section avatar>
            <q-icon :name="link.icon" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ link.label }}</q-item-label>
            <q-item-label caption>{{ link.caption }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer elevated class="bg-grey-8 text-white">
      <q-tabs align="justify" class="text-white">
        <q-route-tab
          v-for="tab in navLinks"
          :key="tab.to"
          :to="{ name: tab.to }"
          :icon="tab.icon"
          :label="tab.label"
        />
      </q-tabs>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface NavLink {
  to: string;
  icon: string;
  label: string;
  caption: string;
}

const navLinks: NavLink[] = [
  {
    to: 'home',
    icon: 'home',
    label: 'Home',
    caption: 'Dashboard',
  },
  {
    to: 'recipes',
    icon: 'menu_book',
    label: 'Recipes',
    caption: 'Browse all recipes',
  },
  {
    to: 'mealplan',
    icon: 'calendar_month',
    label: 'Meal Plan',
    caption: 'Plan your meals',
  },
  {
    to: 'search',
    icon: 'search',
    label: 'Search',
    caption: 'Find recipes',
  },
];

const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>

<style scoped>
.text-decoration-none {
  text-decoration: none;
}
</style>
