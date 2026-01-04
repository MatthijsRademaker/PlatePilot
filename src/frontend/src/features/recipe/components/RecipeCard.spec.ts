import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import RecipeCard from './RecipeCard.vue';
import type { Recipe } from '@features/recipe/types/recipe';

const createMockRecipe = (overrides: Partial<Recipe> = {}): Recipe => ({
  id: '1',
  name: 'Test Recipe',
  description: 'A delicious test recipe',
  ingredients: [{ id: '1', name: 'Salt' }],
  cuisines: [{ id: '1', name: 'Italian' }],
  allergies: [],
  preparationTime: 15,
  cookingTime: 30,
  servings: 4,
  instructions: ['Step 1', 'Step 2'],
  ...overrides,
});

describe('RecipeCard', () => {
  it('renders recipe name', () => {
    const recipe = createMockRecipe({ name: 'Spaghetti Carbonara' });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('Spaghetti Carbonara');
  });

  it('renders recipe description', () => {
    const recipe = createMockRecipe({ description: 'Creamy Italian pasta' });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('Creamy Italian pasta');
  });

  it('computes total time correctly', () => {
    const recipe = createMockRecipe({
      preparationTime: 10,
      cookingTime: 25,
    });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    // 10 + 25 = 35 min
    expect(wrapper.text()).toContain('35 min');
  });

  it('emits click event with recipe when clicked', async () => {
    const recipe = createMockRecipe();

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    await wrapper.trigger('click');

    expect(wrapper.emitted('click')).toBeTruthy();
    expect(wrapper.emitted('click')![0]).toEqual([recipe]);
  });

  it('displays servings', () => {
    const recipe = createMockRecipe({ servings: 6 });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('6 servings');
  });
});
