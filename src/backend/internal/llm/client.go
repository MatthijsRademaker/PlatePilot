package llm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"

	"github.com/platepilot/backend/internal/common/config"
)

// Client wraps the OpenAI SDK for use with Ollama (dev) or OpenAI (prod)
type Client struct {
	client      openai.Client
	model       string
	embedModel  string
	maxTokens   int
	temperature float32
	timeout     time.Duration
	logger      *slog.Logger
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// ChatResponse represents the response from a chat completion
type ChatResponse struct {
	Content      string
	TokensUsed   int
	FinishReason string
}

// EmbeddingResponse represents the response from an embedding request
type EmbeddingResponse struct {
	Embedding  []float32
	TokensUsed int
}

// NewClient creates a new LLM client from configuration
func NewClient(cfg config.LLM, logger *slog.Logger) (*Client, error) {
	if !cfg.IsConfigured() {
		return nil, fmt.Errorf("LLM not configured: missing base_url or api_key")
	}

	// Create OpenAI client with custom base URL for Ollama compatibility
	client := openai.NewClient(
		option.WithBaseURL(cfg.BaseURL),
		option.WithAPIKey(cfg.APIKey),
	)

	return &Client{
		client:      client,
		model:       cfg.Model,
		embedModel:  cfg.EmbedModel,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		timeout:     cfg.Timeout,
		logger:      logger,
	}, nil
}

// Chat sends a chat completion request and returns the response
func (c *Client) Chat(ctx context.Context, messages []ChatMessage) (*ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Convert messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case "system":
			openaiMessages[i] = openai.SystemMessage(msg.Content)
		case "user":
			openaiMessages[i] = openai.UserMessage(msg.Content)
		case "assistant":
			openaiMessages[i] = openai.AssistantMessage(msg.Content)
		default:
			openaiMessages[i] = openai.UserMessage(msg.Content)
		}
	}

	c.logger.Debug("sending chat request",
		"model", c.model,
		"message_count", len(messages),
	)

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       c.model,
		Messages:    openaiMessages,
		MaxTokens:   openai.Int(int64(c.maxTokens)),
		Temperature: openai.Float(float64(c.temperature)),
	})
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	choice := resp.Choices[0]
	return &ChatResponse{
		Content:      choice.Message.Content,
		TokensUsed:   int(resp.Usage.TotalTokens),
		FinishReason: string(choice.FinishReason),
	}, nil
}

// ChatWithJSON sends a chat completion request expecting JSON response
func (c *Client) ChatWithJSON(ctx context.Context, messages []ChatMessage) (*ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Convert messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case "system":
			openaiMessages[i] = openai.SystemMessage(msg.Content)
		case "user":
			openaiMessages[i] = openai.UserMessage(msg.Content)
		case "assistant":
			openaiMessages[i] = openai.AssistantMessage(msg.Content)
		default:
			openaiMessages[i] = openai.UserMessage(msg.Content)
		}
	}

	c.logger.Debug("sending JSON chat request",
		"model", c.model,
		"message_count", len(messages),
	)

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       c.model,
		Messages:    openaiMessages,
		MaxTokens:   openai.Int(int64(c.maxTokens)),
		Temperature: openai.Float(float64(c.temperature)),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: ptr(shared.NewResponseFormatJSONObjectParam()),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	choice := resp.Choices[0]
	return &ChatResponse{
		Content:      choice.Message.Content,
		TokensUsed:   int(resp.Usage.TotalTokens),
		FinishReason: string(choice.FinishReason),
	}, nil
}

// Embed generates embeddings for the given text
func (c *Client) Embed(ctx context.Context, text string) (*EmbeddingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.logger.Debug("generating embedding",
		"model", c.embedModel,
		"text_length", len(text),
	)

	resp, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: c.embedModel,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("embedding generation failed: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Convert float64 to float32
	embedding := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float32(v)
	}

	return &EmbeddingResponse{
		Embedding:  embedding,
		TokensUsed: int(resp.Usage.TotalTokens),
	}, nil
}

// EmbedBatch generates embeddings for multiple texts
func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.logger.Debug("generating batch embeddings",
		"model", c.embedModel,
		"batch_size", len(texts),
	)

	resp, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: c.embedModel,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("batch embedding generation failed: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embedding := make([]float32, len(data.Embedding))
		for j, v := range data.Embedding {
			embedding[j] = float32(v)
		}
		embeddings[i] = embedding
	}

	return embeddings, nil
}

// Model returns the configured chat model name
func (c *Client) Model() string {
	return c.model
}

// EmbedModel returns the configured embedding model name
func (c *Client) EmbedModel() string {
	return c.embedModel
}

// HealthCheck verifies the LLM service is reachable
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple models list to verify connectivity
	_, err := c.client.Models.List(ctx)
	if err != nil {
		return fmt.Errorf("LLM health check failed: %w", err)
	}
	return nil
}

// ptr returns a pointer to the given value
func ptr[T any](v T) *T {
	return &v
}
