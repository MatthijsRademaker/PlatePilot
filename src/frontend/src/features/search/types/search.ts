export interface SearchFilters {
  query: string;
  cuisineIds?: string[];
  excludeAllergyIds?: string[];
  maxPrepTime?: number;
  maxCookTime?: number;
}
