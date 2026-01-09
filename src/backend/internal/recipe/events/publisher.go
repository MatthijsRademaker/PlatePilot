package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/common/dto"
	"github.com/platepilot/backend/internal/common/events"
)

// Publisher publishes recipe events to RabbitMQ
type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   *slog.Logger
}

// PublisherConfig contains configuration for the event publisher
type PublisherConfig struct {
	URL          string
	ExchangeName string
}

// NewPublisher creates a new RabbitMQ event publisher
func NewPublisher(cfg PublisherConfig, logger *slog.Logger) (*Publisher, error) {
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

	logger.Info("connected to RabbitMQ", "exchange", cfg.ExchangeName)

	return &Publisher{
		conn:     conn,
		channel:  ch,
		exchange: cfg.ExchangeName,
		logger:   logger,
	}, nil
}

// Publish publishes an event to RabbitMQ
func (p *Publisher) Publish(ctx context.Context, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// Determine routing key from event type
	routingKey := routingKeyForEvent(event)

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	p.logger.Debug("published event",
		"eventType", event.EventType(),
		"eventId", event.EventID(),
		"aggregateId", event.AggregateID(),
		"routingKey", routingKey,
	)

	return nil
}

// Close closes the connection
func (p *Publisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}

// PublishRecipeUpserted publishes a RecipeUpsertedEvent.
func (p *Publisher) PublishRecipeUpserted(ctx context.Context, recipe *domain.Recipe) error {
	recipeDTO := dto.FromRecipe(recipe)
	event := events.NewRecipeUpsertedEvent(recipeDTO)

	p.logger.Info("publishing recipe upserted event",
		"recipeId", recipe.ID,
		"recipeName", recipe.Name,
	)

	return p.Publish(ctx, event)
}

// PublishRecipeDeleted publishes a RecipeDeletedEvent.
func (p *Publisher) PublishRecipeDeleted(ctx context.Context, recipeID, userID uuid.UUID) error {
	event := events.NewRecipeDeletedEvent(recipeID, userID)

	p.logger.Info("publishing recipe deleted event",
		"recipeId", recipeID,
	)

	return p.Publish(ctx, event)
}

// routingKeyForEvent returns the routing key for a given event
func routingKeyForEvent(event events.Event) string {
	switch event.EventType() {
	case "RecipeUpsertedEvent":
		return "recipe.upserted"
	case "RecipeDeletedEvent":
		return "recipe.deleted"
	default:
		return "recipe.unknown"
	}
}
