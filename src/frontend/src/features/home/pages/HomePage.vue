<template>
  <q-page padding>
    <div class="text-center q-py-xl">
      <q-icon name="restaurant_menu" size="80px" color="primary" />
      <h1 class="text-h3 q-mt-md q-mb-sm">PlatePilot</h1>
      <p class="text-subtitle1 text-grey">Your intelligent meal planning companion</p>
    </div>

    <!-- Today's Plan Card -->
    <div class="row justify-center q-mb-lg">
      <div class="col-12 col-md-8 col-lg-6">
        <TodayPlanCard />
      </div>
    </div>

    <div class="row q-col-gutter-lg justify-center q-mt-lg">
      <div class="col-12 col-sm-6 col-md-4">
        <q-card class="cursor-pointer" @click="$router.push({ name: 'recipes' })">
          <q-card-section class="text-center">
            <q-icon name="menu_book" size="48px" color="primary" />
            <div class="text-h6 q-mt-sm">Browse Recipes</div>
            <div class="text-caption text-grey">
              Explore our collection of delicious recipes
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-4">
        <q-card class="cursor-pointer" @click="$router.push({ name: 'mealplan' })">
          <q-card-section class="text-center">
            <q-icon name="calendar_month" size="48px" color="secondary" />
            <div class="text-h6 q-mt-sm">Meal Plan</div>
            <div class="text-caption text-grey">
              Plan your meals for the week
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-4">
        <q-card class="cursor-pointer" @click="$router.push({ name: 'search' })">
          <q-card-section class="text-center">
            <q-icon name="search" size="48px" color="accent" />
            <div class="text-h6 q-mt-sm">Search</div>
            <div class="text-caption text-grey">
              Find the perfect recipe
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { TodayPlanCard } from '../components';
</script>
