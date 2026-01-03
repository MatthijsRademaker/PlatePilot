package llm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/common/domain"
)

// EmbeddingGenerator implements vector.Generator using LLM embeddings
type EmbeddingGenerator struct {
	client     *Client
	dimensions int
	logger     *slog.Logger
}

// NewEmbeddingGenerator creates a new LLM-based embedding generator
func NewEmbeddingGenerator(client *Client, dimensions int, logger *slog.Logger) *EmbeddingGenerator {
	return &EmbeddingGenerator{
		client:     client,
		dimensions: dimensions,
		logger:     logger,
	}
}

// Generate creates a vector embedding from text using the LLM embedding API
func (g *EmbeddingGenerator) Generate(text string) pgvector.Vector {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := g.client.Embed(ctx, text)
	if err != nil {
		g.logger.Error("failed to generate embedding", "error", err, "text_length", len(text))
		// Return zero vector on error
		return pgvector.NewVector(make([]float32, g.dimensions))
	}

	// Handle dimension mismatch
	embedding := resp.Embedding
	if len(embedding) != g.dimensions {
		g.logger.Warn("embedding dimension mismatch",
			"expected", g.dimensions,
			"got", len(embedding),
		)
		embedding = adjustDimensions(embedding, g.dimensions)
	}

	return pgvector.NewVector(embedding)
}

// GenerateForRecipe creates a vector embedding for a recipe
func (g *EmbeddingGenerator) GenerateForRecipe(recipe *domain.Recipe) pgvector.Vector {
	// Build rich text representation of the recipe
	text := buildRecipeText(recipe)
	return g.Generate(text)
}

// GenerateWithContext creates an embedding with explicit context
func (g *EmbeddingGenerator) GenerateWithContext(ctx context.Context, text string) (pgvector.Vector, error) {
	resp, err := g.client.Embed(ctx, text)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("generate embedding: %w", err)
	}

	embedding := resp.Embedding
	if len(embedding) != g.dimensions {
		embedding = adjustDimensions(embedding, g.dimensions)
	}

	return pgvector.NewVector(embedding), nil
}

// GenerateBatch creates embeddings for multiple texts
func (g *EmbeddingGenerator) GenerateBatch(ctx context.Context, texts []string) ([]pgvector.Vector, error) {
	embeddings, err := g.client.EmbedBatch(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("generate batch embeddings: %w", err)
	}

	vectors := make([]pgvector.Vector, len(embeddings))
	for i, embedding := range embeddings {
		if len(embedding) != g.dimensions {
			embedding = adjustDimensions(embedding, g.dimensions)
		}
		vectors[i] = pgvector.NewVector(embedding)
	}

	return vectors, nil
}

// Dimensions returns the configured embedding dimensions
func (g *EmbeddingGenerator) Dimensions() int {
	return g.dimensions
}

// buildRecipeText creates a rich text representation for embedding
func buildRecipeText(recipe *domain.Recipe) string {
	var parts []string

	// Recipe name (important)
	parts = append(parts, recipe.Name)

	// Description
	if recipe.Description != "" {
		parts = append(parts, recipe.Description)
	}

	// Main ingredient
	if recipe.MainIngredient != nil {
		parts = append(parts, "Main ingredient: "+recipe.MainIngredient.Name)
	}

	// Cuisine
	if recipe.Cuisine != nil {
		parts = append(parts, "Cuisine: "+recipe.Cuisine.Name)
	}

	// Other ingredients
	if len(recipe.Ingredients) > 0 {
		var ingredientNames []string
		for _, ing := range recipe.Ingredients {
			ingredientNames = append(ingredientNames, ing.Name)
		}
		parts = append(parts, "Ingredients: "+joinStrings(ingredientNames, ", "))
	}

	// Tags from metadata
	if len(recipe.Metadata.Tags) > 0 {
		parts = append(parts, "Tags: "+joinStrings(recipe.Metadata.Tags, ", "))
	}

	return joinStrings(parts, ". ")
}

// adjustDimensions pads or truncates embedding to target dimensions
func adjustDimensions(embedding []float32, target int) []float32 {
	if len(embedding) == target {
		return embedding
	}

	result := make([]float32, target)
	if len(embedding) > target {
		// Truncate
		copy(result, embedding[:target])
	} else {
		// Pad with zeros
		copy(result, embedding)
	}
	return result
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
