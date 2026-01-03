package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/platepilot/backend/internal/llm"
)

// Agent defines the interface for all AI agents
type Agent interface {
	// Name returns the agent's identifier
	Name() string

	// Execute runs the agent with the given input and returns structured output
	Execute(ctx context.Context, input AgentInput) (*AgentOutput, error)
}

// AgentInput contains the input for an agent execution
type AgentInput struct {
	// UserMessage is the primary user query or instruction
	UserMessage string

	// Context provides additional structured context for the agent
	Context map[string]any
}

// AgentOutput contains the result of an agent execution
type AgentOutput struct {
	// RawContent is the raw response from the LLM
	RawContent string

	// Structured contains the parsed JSON response (if applicable)
	Structured any

	// TokensUsed is the total tokens consumed
	TokensUsed int

	// Cached indicates if this response came from cache
	Cached bool
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	name         string
	systemPrompt string
	client       *llm.Client
	cache        Cache
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(name, systemPrompt string, client *llm.Client, cache Cache) *BaseAgent {
	return &BaseAgent{
		name:         name,
		systemPrompt: systemPrompt,
		client:       client,
		cache:        cache,
	}
}

// Name returns the agent's name
func (a *BaseAgent) Name() string {
	return a.name
}

// ExecuteWithJSON sends a request expecting a JSON response
func (a *BaseAgent) ExecuteWithJSON(ctx context.Context, input AgentInput, result any) (*AgentOutput, error) {
	// Check cache first
	cacheKey := a.cacheKey(input)
	if cached, ok := a.cache.Get(cacheKey); ok {
		if err := json.Unmarshal(cached, result); err == nil {
			return &AgentOutput{
				RawContent: string(cached),
				Structured: result,
				Cached:     true,
			}, nil
		}
	}

	// Build messages
	messages := []llm.ChatMessage{
		{Role: "system", Content: a.systemPrompt},
		{Role: "user", Content: a.buildUserMessage(input)},
	}

	// Call LLM
	resp, err := a.client.ChatWithJSON(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("agent %s failed: %w", a.name, err)
	}

	// Parse response
	if err := json.Unmarshal([]byte(resp.Content), result); err != nil {
		return nil, fmt.Errorf("agent %s: failed to parse response: %w", a.name, err)
	}

	// Cache successful response
	a.cache.Set(cacheKey, []byte(resp.Content), DefaultCacheTTL)

	return &AgentOutput{
		RawContent: resp.Content,
		Structured: result,
		TokensUsed: resp.TokensUsed,
		Cached:     false,
	}, nil
}

// ExecuteWithText sends a request expecting a text response
func (a *BaseAgent) ExecuteWithText(ctx context.Context, input AgentInput) (*AgentOutput, error) {
	// Check cache first
	cacheKey := a.cacheKey(input)
	if cached, ok := a.cache.Get(cacheKey); ok {
		return &AgentOutput{
			RawContent: string(cached),
			Cached:     true,
		}, nil
	}

	// Build messages
	messages := []llm.ChatMessage{
		{Role: "system", Content: a.systemPrompt},
		{Role: "user", Content: a.buildUserMessage(input)},
	}

	// Call LLM
	resp, err := a.client.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("agent %s failed: %w", a.name, err)
	}

	// Cache successful response
	a.cache.Set(cacheKey, []byte(resp.Content), DefaultCacheTTL)

	return &AgentOutput{
		RawContent: resp.Content,
		TokensUsed: resp.TokensUsed,
		Cached:     false,
	}, nil
}

// buildUserMessage constructs the user message from input
func (a *BaseAgent) buildUserMessage(input AgentInput) string {
	if len(input.Context) == 0 {
		return input.UserMessage
	}

	// Include context as JSON in the message
	contextJSON, err := json.MarshalIndent(input.Context, "", "  ")
	if err != nil {
		return input.UserMessage
	}

	return fmt.Sprintf("Context:\n```json\n%s\n```\n\nRequest: %s", string(contextJSON), input.UserMessage)
}

// cacheKey generates a cache key for the given input
func (a *BaseAgent) cacheKey(input AgentInput) string {
	// Simple hash of agent name + input
	data, _ := json.Marshal(input)
	return fmt.Sprintf("%s:%x", a.name, hashBytes(data))
}

// hashBytes returns a simple hash of the given bytes
func hashBytes(data []byte) uint64 {
	var hash uint64 = 14695981039346656037 // FNV-1a offset basis
	for _, b := range data {
		hash ^= uint64(b)
		hash *= 1099511628211 // FNV-1a prime
	}
	return hash
}

// RecipeSuggestionOutput is the structured output for recipe suggestions
type RecipeSuggestionOutput struct {
	Recipes []SuggestedRecipe `json:"recipes"`
}

// SuggestedRecipe represents a single recipe suggestion
type SuggestedRecipe struct {
	RecipeID  string  `json:"recipe_id"`
	Score     float64 `json:"score"`
	Reasoning string  `json:"reasoning"`
}

// MealPlanOutput is the structured output for meal planning
type MealPlanOutput struct {
	Days []MealPlanDay `json:"days"`
}

// MealPlanDay represents a single day in a meal plan
type MealPlanDay struct {
	Day   string `json:"day"`
	Meals []Meal `json:"meals"`
}

// Meal represents a single meal in a day
type Meal struct {
	MealType string `json:"meal_type"` // "breakfast", "lunch", "dinner", "snack"
	RecipeID string `json:"recipe_id"`
	Notes    string `json:"notes,omitempty"`
}
