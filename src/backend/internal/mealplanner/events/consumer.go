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
	case "RecipeUpsertedEvent":
		return c.handleRecipeUpserted(ctx, msg.Body)
	case "RecipeDeletedEvent":
		return c.handleRecipeDeleted(ctx, msg.Body)
	default:
		c.logger.Warn("unknown event type", "type", envelope.Type)
		return nil // Acknowledge unknown events to prevent redelivery
	}
}

func (c *Consumer) handleRecipeUpserted(ctx context.Context, body []byte) error {
	var event RecipeUpsertedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal recipe upserted event: %w", err)
	}

	c.logger.Info("handling recipe upserted event",
		"eventId", event.ID,
		"recipeId", event.Recipe.ID,
		"recipeName", event.Recipe.Name,
	)

	// Convert DTO to repository model and upsert
	recipe := event.Recipe.ToRepositoryModel()
	lines := event.Recipe.ToIngredientLineModels()
	if err := c.repo.Upsert(ctx, recipe, lines); err != nil {
		return fmt.Errorf("upsert recipe: %w", err)
	}

	c.logger.Info("recipe upserted to read model", "recipeId", event.Recipe.ID)
	return nil
}

func (c *Consumer) handleRecipeDeleted(ctx context.Context, body []byte) error {
	var event RecipeDeletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal recipe deleted event: %w", err)
	}

	c.logger.Info("handling recipe deleted event",
		"eventId", event.ID,
		"aggregateId", event.AggregateId,
	)

	if err := c.repo.Delete(ctx, event.AggregateId); err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}

	c.logger.Info("recipe deleted from read model", "recipeId", event.AggregateId)
	return nil
}

// EventEnvelope is the common structure for all events
type EventEnvelope struct {
	ID               uuid.UUID `json:"id"`
	Type             string    `json:"type"`
	OccurredOn        string    `json:"occurredOn"`
	AggregateId       uuid.UUID `json:"aggregateId"`
	SchemaVersion     int       `json:"schemaVersion"`
	AggregateVersion  int       `json:"aggregateVersion"`
}

// RecipeUpsertedEvent represents a recipe upsert event
type RecipeUpsertedEvent struct {
	EventEnvelope
	Recipe RecipeDTO `json:"recipe"`
}

// RecipeDeletedEvent represents a recipe delete event
type RecipeDeletedEvent struct {
	EventEnvelope
	UserID    uuid.UUID `json:"userId"`
	DeletedAt string    `json:"deletedAt"`
}

// RecipeDTO is the recipe data in events
type RecipeDTO struct {
	ID               uuid.UUID          `json:"id"`
	UserID           uuid.UUID          `json:"userId"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	PrepTimeMinutes  int                `json:"prepTimeMinutes"`
	CookTimeMinutes  int                `json:"cookTimeMinutes"`
	TotalTimeMinutes int                `json:"totalTimeMinutes"`
	Servings         int                `json:"servings"`
	YieldQuantity    *float64           `json:"yieldQuantity"`
	YieldUnit        string             `json:"yieldUnit"`
	MainIngredient   IngredientDTO      `json:"mainIngredient"`
	Cuisine          CuisineDTO         `json:"cuisine"`
	IngredientLines  []IngredientLineDTO `json:"ingredientLines"`
	Allergies        []AllergyDTO       `json:"allergies"`
	Tags             []string           `json:"tags"`
	ImageURL         string             `json:"imageUrl"`
	Nutrition        RecipeNutritionDTO `json:"nutrition"`
	SearchVector     []float32          `json:"searchVector"`
}

// IngredientDTO represents an ingredient in events
type IngredientDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// IngredientLineDTO represents an ingredient line item in events
type IngredientLineDTO struct {
	Ingredient    IngredientDTO `json:"ingredient"`
	QuantityValue *float64      `json:"quantityValue"`
	QuantityText  string        `json:"quantityText"`
	Unit          string        `json:"unit"`
	IsOptional    bool          `json:"isOptional"`
	Note          string        `json:"note"`
	SortOrder     int           `json:"sortOrder"`
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

// RecipeNutritionDTO represents aggregated nutrition in events
type RecipeNutritionDTO struct {
	CaloriesTotal      int     `json:"caloriesTotal"`
	CaloriesPerServing int     `json:"caloriesPerServing"`
	ProteinG           float64 `json:"proteinG"`
	CarbsG             float64 `json:"carbsG"`
	FatG               float64 `json:"fatG"`
	FiberG             float64 `json:"fiberG"`
	SugarG             float64 `json:"sugarG"`
	SodiumMg           float64 `json:"sodiumMg"`
}

// ToRepositoryModel converts the DTO to a repository model
func (d *RecipeDTO) ToRepositoryModel() *repository.Recipe {
	tags := d.Tags
	if tags == nil {
		tags = []string{}
	}

	ingredientIDs := make([]uuid.UUID, len(d.IngredientLines))
	for i, line := range d.IngredientLines {
		ingredientIDs[i] = line.Ingredient.ID
	}

	allergyIDs := make([]uuid.UUID, len(d.Allergies))
	for i, a := range d.Allergies {
		allergyIDs[i] = a.ID
	}

	return &repository.Recipe{
		ID:                 d.ID,
		UserID:             d.UserID,
		Name:               d.Name,
		Description:        d.Description,
		PrepTimeMinutes:    d.PrepTimeMinutes,
		CookTimeMinutes:    d.CookTimeMinutes,
		TotalTimeMinutes:   d.TotalTimeMinutes,
		Servings:           d.Servings,
		YieldQuantity:      d.YieldQuantity,
		YieldUnit:          d.YieldUnit,
		SearchVector:       vectorFromSlice(d.SearchVector),
		CuisineID:          d.Cuisine.ID,
		CuisineName:        d.Cuisine.Name,
		MainIngredientID:   d.MainIngredient.ID,
		MainIngredientName: d.MainIngredient.Name,
		IngredientIDs:      ingredientIDs,
		AllergyIDs:         allergyIDs,
		ImageURL:           d.ImageURL,
		Tags:               tags,
		CaloriesTotal:      d.Nutrition.CaloriesTotal,
		CaloriesPerServing: d.Nutrition.CaloriesPerServing,
		ProteinG:           d.Nutrition.ProteinG,
		CarbsG:             d.Nutrition.CarbsG,
		FatG:               d.Nutrition.FatG,
		FiberG:             d.Nutrition.FiberG,
		SugarG:             d.Nutrition.SugarG,
		SodiumMg:           d.Nutrition.SodiumMg,
	}
}

// ToIngredientLineModels converts ingredient line DTOs to repository models.
func (d *RecipeDTO) ToIngredientLineModels() []repository.IngredientLine {
	lines := make([]repository.IngredientLine, 0, len(d.IngredientLines))
	for _, line := range d.IngredientLines {
		lines = append(lines, repository.IngredientLine{
			RecipeID:       d.ID,
			IngredientID:   line.Ingredient.ID,
			IngredientName: line.Ingredient.Name,
			QuantityValue:  line.QuantityValue,
			QuantityText:   line.QuantityText,
			Unit:           line.Unit,
			IsOptional:     line.IsOptional,
			Note:           line.Note,
			SortOrder:      line.SortOrder,
		})
	}
	return lines
}

func vectorFromSlice(v []float32) pgvector.Vector {
	return pgvector.NewVector(v)
}
