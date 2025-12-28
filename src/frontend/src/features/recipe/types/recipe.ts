export interface Ingredient {
  id: string;
  name: string;
}

export interface Cuisine {
  id: string;
  name: string;
}

export interface Allergy {
  id: string;
  name: string;
}

export interface Recipe {
  id: string;
  name: string;
  description: string;
  ingredients: Ingredient[];
  cuisines: Cuisine[];
  allergies: Allergy[];
  preparationTime: number;
  cookingTime: number;
  servings: number;
  instructions: string[];
  imageUrl?: string;
}

export interface CreateRecipeRequest {
  name: string;
  description: string;
  ingredientIds: string[];
  cuisineIds: string[];
  allergyIds: string[];
  preparationTime: number;
  cookingTime: number;
  servings: number;
  instructions: string[];
  imageUrl?: string;
}
