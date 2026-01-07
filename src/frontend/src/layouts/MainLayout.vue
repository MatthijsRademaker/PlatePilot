<template>
  <q-layout view="hHh Lpr fFf">
    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="main-footer">
      <nav class="footer-nav">
        <router-link
          v-for="tab in navLinks"
          :key="tab.to"
          :to="{ name: tab.to }"
          class="nav-item"
          :class="{ 'nav-item--active': isActive(tab.to) }"
        >
          <div class="nav-icon">
            <q-icon :name="tab.icon" size="22px" />
          </div>
          <span class="nav-label">{{ tab.label }}</span>
        </router-link>
      </nav>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

interface NavLink {
  to: string;
  icon: string;
  label: string;
}

const route = useRoute();

const navLinks: NavLink[] = [
  {
    to: 'home',
    icon: 'mdi-home-variant',
    label: 'Home',
  },
  {
    to: 'recipes',
    icon: 'mdi-book-open-page-variant',
    label: 'Recipes',
  },
  {
    to: 'mealplan',
    icon: 'mdi-calendar-month',
    label: 'Plan',
  },
  {
    to: 'search',
    icon: 'mdi-magnify',
    label: 'Search',
  },
];

function isActive(routeName: string): boolean {
  return route.name === routeName;
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&display=swap');

.main-footer {
  background: #2d1f1a;
  padding-bottom: env(safe-area-inset-bottom);
}

.footer-nav {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 8px 0;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  text-decoration: none;
  transition: all 0.2s ease;
  border-radius: 12px;

  .nav-icon {
    color: rgba(255, 255, 255, 0.5);
    transition: all 0.2s ease;
  }

  .nav-label {
    font-family: 'DM Sans', sans-serif;
    font-size: 11px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.5);
    transition: all 0.2s ease;
  }

  &:active {
    transform: scale(0.95);
  }

  &--active {
    .nav-icon {
      color: #ff7f50;
    }

    .nav-label {
      color: #ff7f50;
      font-weight: 600;
    }
  }

  &:not(.nav-item--active):hover {
    .nav-icon,
    .nav-label {
      color: rgba(255, 255, 255, 0.8);
    }
  }
}
</style>
