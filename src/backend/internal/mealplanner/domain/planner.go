package domain

import (
	"context"
	"math"
	"sort"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// SuggestionRequest contains the parameters for suggesting recipes
type SuggestionRequest struct {
	UserID                 uuid.UUID
	DailyConstraints       []DailyConstraints
	AlreadySelectedRecipes []uuid.UUID
	Amount                 int
}

// DailyConstraints represents constraints for a single day's meal
type DailyConstraints struct {
	IngredientConstraints []uuid.UUID
	CuisineConstraints    []uuid.UUID
}

// Planner suggests recipes based on constraints and diversity
type Planner struct {
	repo RecipeRepository
}

// NewPlanner creates a new meal planner
func NewPlanner(repo RecipeRepository) *Planner {
	return &Planner{repo: repo}
}

// SuggestMeals suggests recipes based on the given request
func (p *Planner) SuggestMeals(ctx context.Context, req SuggestionRequest) ([]uuid.UUID, error) {
	// Get all recipes from the repository
	recipes, err := p.repo.GetAll(ctx, req.UserID, 1000, 0) // Get up to 1000 recipes
	if err != nil {
		return nil, err
	}

	// Filter by constraints
	filtered := p.filterByConstraints(recipes, req.DailyConstraints)

	// Remove already selected recipes
	filtered = p.removeSelected(filtered, req.AlreadySelectedRecipes)

	if len(filtered) == 0 {
		return []uuid.UUID{}, nil
	}

	// Score recipes for diversity
	scored := p.scoreForDiversity(filtered, req.AlreadySelectedRecipes, recipes)

	// Sort by score (higher is better)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Take top N
	amount := req.Amount
	if amount > len(scored) {
		amount = len(scored)
	}

	result := make([]uuid.UUID, amount)
	for i := 0; i < amount; i++ {
		result[i] = scored[i].id
	}

	return result, nil
}

type scoredRecipe struct {
	id    uuid.UUID
	score float64
}

func (p *Planner) filterByConstraints(recipes []Recipe, constraints []DailyConstraints) []Recipe {
	if len(constraints) == 0 {
		return recipes
	}

	var filtered []Recipe
	for _, recipe := range recipes {
		if p.matchesAnyConstraint(recipe, constraints) {
			filtered = append(filtered, recipe)
		}
	}
	return filtered
}

func (p *Planner) matchesAnyConstraint(recipe Recipe, constraints []DailyConstraints) bool {
	// If there are no constraints, all recipes match
	if len(constraints) == 0 {
		return true
	}

	// Recipe must match at least one day's constraints
	for _, daily := range constraints {
		if p.matchesDailyConstraint(recipe, daily) {
			return true
		}
	}
	return false
}

func (p *Planner) matchesDailyConstraint(recipe Recipe, constraint DailyConstraints) bool {
	// Check cuisine constraints (if any specified, must match one)
	if len(constraint.CuisineConstraints) > 0 {
		cuisineMatch := false
		for _, cuisineID := range constraint.CuisineConstraints {
			if recipe.CuisineID == cuisineID {
				cuisineMatch = true
				break
			}
		}
		if !cuisineMatch {
			return false
		}
	}

	// Check ingredient constraints (if any specified, must contain one)
	if len(constraint.IngredientConstraints) > 0 {
		ingredientMatch := false
		for _, ingredientID := range constraint.IngredientConstraints {
			if recipe.MainIngredientID == ingredientID {
				ingredientMatch = true
				break
			}
			for _, recipeIngID := range recipe.IngredientIDs {
				if recipeIngID == ingredientID {
					ingredientMatch = true
					break
				}
			}
			if ingredientMatch {
				break
			}
		}
		if !ingredientMatch {
			return false
		}
	}

	return true
}

func (p *Planner) removeSelected(recipes []Recipe, selected []uuid.UUID) []Recipe {
	if len(selected) == 0 {
		return recipes
	}

	selectedSet := make(map[uuid.UUID]bool)
	for _, id := range selected {
		selectedSet[id] = true
	}

	var filtered []Recipe
	for _, recipe := range recipes {
		if !selectedSet[recipe.ID] {
			filtered = append(filtered, recipe)
		}
	}
	return filtered
}

func (p *Planner) scoreForDiversity(candidates []Recipe, selected []uuid.UUID, allRecipes []Recipe) []scoredRecipe {
	// Build a map of selected recipe vectors for diversity calculation
	selectedVectors := make([]pgvector.Vector, 0, len(selected))
	recipeMap := make(map[uuid.UUID]Recipe)
	for _, r := range allRecipes {
		recipeMap[r.ID] = r
	}

	for _, id := range selected {
		if r, ok := recipeMap[id]; ok {
			selectedVectors = append(selectedVectors, r.SearchVector)
		}
	}

	scored := make([]scoredRecipe, len(candidates))
	for i, candidate := range candidates {
		diversityScore := p.calculateDiversityScore(candidate.SearchVector, selectedVectors)
		scored[i] = scoredRecipe{
			id:    candidate.ID,
			score: diversityScore,
		}
	}

	return scored
}

func (p *Planner) calculateDiversityScore(candidate pgvector.Vector, selected []pgvector.Vector) float64 {
	if len(selected) == 0 {
		return 1.0 // Maximum diversity when nothing selected
	}

	var totalSimilarity float64
	for _, s := range selected {
		similarity := cosineSimilarity(candidate.Slice(), s.Slice())
		totalSimilarity += similarity
	}

	avgSimilarity := totalSimilarity / float64(len(selected))
	// Return diversity score (1 - similarity, so higher is more diverse)
	return 1.0 - avgSimilarity
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
