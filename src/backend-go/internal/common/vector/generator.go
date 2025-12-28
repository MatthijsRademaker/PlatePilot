package vector

import (
	"hash/fnv"
	"math"
	"strings"

	"github.com/pgvector/pgvector-go"
	"github.com/platepilot/backend/internal/common/domain"
)

const VectorDimensions = 128

// Generator creates vector embeddings for recipes
type Generator interface {
	Generate(text string) pgvector.Vector
	GenerateForRecipe(recipe *domain.Recipe) pgvector.Vector
}

// HashGenerator is a POC implementation using hash-based vectors.
// This should be replaced with Azure OpenAI embeddings in production.
type HashGenerator struct{}

// NewHashGenerator creates a new hash-based vector generator
func NewHashGenerator() *HashGenerator {
	return &HashGenerator{}
}

// Generate creates a vector embedding from text using hash-based approach
func (g *HashGenerator) Generate(text string) pgvector.Vector {
	words := strings.Fields(strings.ToLower(text))
	vector := make([]float32, VectorDimensions)

	for _, word := range words {
		h := fnv.New32a()
		h.Write([]byte(word))
		idx := h.Sum32() % uint32(VectorDimensions)
		vector[idx] += 1.0
	}

	// Normalize the vector
	normalize(vector)

	return pgvector.NewVector(vector)
}

// GenerateForRecipe creates a vector embedding for a recipe
func (g *HashGenerator) GenerateForRecipe(recipe *domain.Recipe) pgvector.Vector {
	var mainIngredientName string
	if recipe.MainIngredient != nil {
		mainIngredientName = recipe.MainIngredient.Name
	}

	combined := recipe.Name + " " + recipe.Description + " " + mainIngredientName
	return g.Generate(combined)
}

// normalize normalizes a vector to unit length
func normalize(v []float32) {
	var sumSquares float64
	for _, val := range v {
		sumSquares += float64(val) * float64(val)
	}

	if sumSquares == 0 {
		return
	}

	norm := float32(1.0 / math.Sqrt(sumSquares))
	for i := range v {
		v[i] *= norm
	}
}
