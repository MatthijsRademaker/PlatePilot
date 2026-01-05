import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import RecipeCard from './RecipeCard.vue';
import type { HandlerRecipeJSON } from '@/api/generated/models';

const createMockRecipe = (overrides: Partial<HandlerRecipeJSON> = {}): HandlerRecipeJSON => ({
  id: '1',
  name: 'Test Recipe',
  description: 'A delicious test recipe',
  ingredients: [{ id: '1', name: 'Salt' }],
  cuisine: { id: '1', name: 'Italian' },
  prepTime: '15 min',
  cookTime: '30 min',
  directions: ['Step 1', 'Step 2'],
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

  it('displays prep time', () => {
    const recipe = createMockRecipe({
      prepTime: '10 min',
    });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('Prep: 10 min');
  });

  it('displays cook time', () => {
    const recipe = createMockRecipe({
      cookTime: '25 min',
    });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('Cook: 25 min');
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

  it('displays cuisine name', () => {
    const recipe = createMockRecipe({ cuisine: { id: '1', name: 'Mexican' } });

    const wrapper = mount(RecipeCard, {
      props: { recipe },
    });

    expect(wrapper.text()).toContain('Mexican');
  });
});
