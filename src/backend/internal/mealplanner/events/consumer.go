package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/platepilot/backend/internal/mealplanner/repository"
)

// Consumer consumes recipe events from RabbitMQ
type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	repo       *repository.Repository
	logger     *slog.Logger
	queueName  string
	exchange   string
	routingKey string
}

// ConsumerConfig contains configuration for the event consumer
type ConsumerConfig struct {
	URL          string
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

// NewConsumer creates a new RabbitMQ event consumer
func NewConsumer(cfg ConsumerConfig, repo *repository.Repository, logger *slog.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.ExchangeName, // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		cfg.QueueName, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		cfg.QueueName,    // queue name
		cfg.RoutingKey,   // routing key
		cfg.ExchangeName, // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("bind queue: %w", err)
	}

	return &Consumer{
		conn:       conn,
		channel:    ch,
		repo:       repo,
		logger:     logger,
		queueName:  cfg.QueueName,
		exchange:   cfg.ExchangeName,
		routingKey: cfg.RoutingKey,
	}, nil
}

// Start starts consuming events
func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer tag
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	c.logger.Info("started consuming events",
		"queue", c.queueName,
		"exchange", c.exchange,
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopping event consumer")
				return
			case msg, ok := <-msgs:
				if !ok {
					c.logger.Warn("message channel closed")
					return
				}
				if err := c.handleMessage(ctx, msg); err != nil {
					c.logger.Error("failed to handle message",
						"error", err,
						"routingKey", msg.RoutingKey,
					)
					// Nack and requeue on error
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Close closes the connection
func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *Consumer) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	c.logger.Debug("received message",
		"routingKey", msg.RoutingKey,
		"contentType", msg.ContentType,
	)

	// Parse the event envelope to determine type
	var envelope EventEnvelope
	if err := json.Unmarshal(msg.Body, &envelope); err != nil {
		return fmt.Errorf("unmarshal envelope: %w", err)
	}

	switch envelope.Type {
	case "RecipeCreatedEvent":
		return c.handleRecipeCreated(ctx, msg.Body)
	case "RecipeUpdatedEvent":
		return c.handleRecipeUpdated(ctx, msg.Body)
	default:
		c.logger.Warn("unknown event type", "type", envelope.Type)
		return nil // Acknowledge unknown events to prevent redelivery
	}
}

func (c *Consumer) handleRecipeCreated(ctx context.Context, body []byte) error {
	var event RecipeCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal recipe created event: %w", err)
	}

	c.logger.Info("handling recipe created event",
		"eventId", event.ID,
		"recipeId", event.Recipe.ID,
		"recipeName", event.Recipe.Name,
	)

	// Convert DTO to repository model and upsert
	recipe := event.Recipe.ToRepositoryModel()
	if err := c.repo.Upsert(ctx, recipe); err != nil {
		return fmt.Errorf("upsert recipe: %w", err)
	}

	c.logger.Info("recipe upserted to read model", "recipeId", event.Recipe.ID)
	return nil
}

func (c *Consumer) handleRecipeUpdated(ctx context.Context, body []byte) error {
	var event RecipeUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal recipe updated event: %w", err)
	}

	c.logger.Info("handling recipe updated event",
		"eventId", event.ID,
		"aggregateId", event.AggregateId,
	)

	// For updates, we'd need the full recipe data
	// This could be fetched from the Recipe API or included in the event
	// For now, just log it
	c.logger.Warn("recipe update handling not fully implemented - would need to fetch from Recipe API")
	return nil
}

// EventEnvelope is the common structure for all events
type EventEnvelope struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	OccurredOn  string    `json:"occurredOn"`
	AggregateId uuid.UUID `json:"aggregateId"`
}

// RecipeCreatedEvent represents a recipe creation event
type RecipeCreatedEvent struct {
	EventEnvelope
	Recipe RecipeDTO `json:"recipe"`
}

// RecipeUpdatedEvent represents a recipe update event
type RecipeUpdatedEvent struct {
	EventEnvelope
}

// RecipeDTO is the recipe data in events
type RecipeDTO struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	PrepTime        string          `json:"prepTime"`
	CookTime        string          `json:"cookTime"`
	MainIngredient  IngredientDTO   `json:"mainIngredient"`
	Cuisine         CuisineDTO      `json:"cuisine"`
	Ingredients     []IngredientDTO `json:"ingredients"`
	Allergies       []AllergyDTO    `json:"allergies"`
	Directions      []string        `json:"directions"`
	NutritionalInfo NutritionalDTO  `json:"nutritionalInfo"`
	Metadata        MetadataDTO     `json:"metadata"`
}

// IngredientDTO represents an ingredient in events
type IngredientDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// CuisineDTO represents a cuisine in events
type CuisineDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// AllergyDTO represents an allergy in events
type AllergyDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NutritionalDTO represents nutritional info in events
type NutritionalDTO struct {
	Calories int `json:"calories"`
}

// MetadataDTO represents metadata in events
type MetadataDTO struct {
	SearchVector  []float32 `json:"searchVector"`
	ImageURL      string    `json:"imageUrl"`
	Tags          []string  `json:"tags"`
	PublishedDate string    `json:"publishedDate"`
}

// ToRepositoryModel converts the DTO to a repository model
func (d *RecipeDTO) ToRepositoryModel() *repository.Recipe {
	ingredientIDs := make([]uuid.UUID, len(d.Ingredients))
	for i, ing := range d.Ingredients {
		ingredientIDs[i] = ing.ID
	}

	allergyIDs := make([]uuid.UUID, len(d.Allergies))
	for i, a := range d.Allergies {
		allergyIDs[i] = a.ID
	}

	return &repository.Recipe{
		ID:                 d.ID,
		Name:               d.Name,
		Description:        d.Description,
		PrepTime:           d.PrepTime,
		CookTime:           d.CookTime,
		SearchVector:       vectorFromSlice(d.Metadata.SearchVector),
		CuisineID:          d.Cuisine.ID,
		CuisineName:        d.Cuisine.Name,
		MainIngredientID:   d.MainIngredient.ID,
		MainIngredientName: d.MainIngredient.Name,
		IngredientIDs:      ingredientIDs,
		AllergyIDs:         allergyIDs,
		Directions:         d.Directions,
		ImageURL:           d.Metadata.ImageURL,
		Tags:               d.Metadata.Tags,
		Calories:           d.NutritionalInfo.Calories,
	}
}

func vectorFromSlice(v []float32) pgvector.Vector {
	return pgvector.NewVector(v)
}
