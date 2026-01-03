package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Environment string         `mapstructure:"environment"`
	LogLevel    string         `mapstructure:"log_level"`
	RecipeAPI   RecipeAPI      `mapstructure:"recipe_api"`
	MealPlanner MealPlannerAPI `mapstructure:"mealplanner_api"`
	BFF         BFFAPI         `mapstructure:"bff"`
	Database    Database       `mapstructure:"database"`
	RabbitMQ    RabbitMQ       `mapstructure:"rabbitmq"`
	LLM         LLM            `mapstructure:"llm"`
}

// RecipeAPI configuration
type RecipeAPI struct {
	HTTPAddress string        `mapstructure:"http_address"`
	GRPCAddress string        `mapstructure:"grpc_address"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// MealPlannerAPI configuration
type MealPlannerAPI struct {
	HTTPAddress string        `mapstructure:"http_address"`
	GRPCAddress string        `mapstructure:"grpc_address"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// BFFAPI configuration
type BFFAPI struct {
	HTTPAddress       string        `mapstructure:"http_address"`
	Timeout           time.Duration `mapstructure:"timeout"`
	RecipeAPIAddress  string        `mapstructure:"recipe_api_address"`
	MealPlanAddress   string        `mapstructure:"mealplan_api_address"`
	CORSAllowedOrigins []string     `mapstructure:"cors_allowed_origins"`
}

// Database configuration
type Database struct {
	RecipeDB      string        `mapstructure:"recipe_db"`
	MealPlannerDB string        `mapstructure:"mealplanner_db"`
	MaxOpenConns  int           `mapstructure:"max_open_conns"`
	MaxIdleConns  int           `mapstructure:"max_idle_conns"`
	ConnMaxLife   time.Duration `mapstructure:"conn_max_life"`
}

// RabbitMQ configuration
type RabbitMQ struct {
	URL          string `mapstructure:"url"`
	ExchangeName string `mapstructure:"exchange_name"`
	ExchangeType string `mapstructure:"exchange_type"`
}

// LLM configuration for OpenAI-compatible APIs (Ollama, OpenAI, Azure)
type LLM struct {
	BaseURL         string        `mapstructure:"base_url"`         // API base URL (Ollama: http://localhost:11434/v1)
	APIKey          string        `mapstructure:"api_key"`          // API key (Ollama: "ollama", OpenAI: real key)
	Model           string        `mapstructure:"model"`            // Chat model (llama3.1, gpt-4o-mini)
	EmbedModel      string        `mapstructure:"embed_model"`      // Embedding model (nomic-embed-text, text-embedding-3-small)
	EmbedDimensions int           `mapstructure:"embed_dimensions"` // Embedding dimensions (768 for nomic, 1536 for OpenAI)
	Timeout         time.Duration `mapstructure:"timeout"`          // Request timeout
	MaxTokens       int           `mapstructure:"max_tokens"`       // Max response tokens
	Temperature     float32       `mapstructure:"temperature"`      // Response temperature (0.0-2.0)
}

// IsConfigured returns true if LLM is properly configured
func (l *LLM) IsConfigured() bool {
	return l.BaseURL != "" && l.APIKey != ""
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	return LoadWithPath("")
}

// LoadWithPath reads configuration from a specific path
func LoadWithPath(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/platepilot")

	// Environment variable settings
	v.SetEnvPrefix("PLATEPILOT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Environment
	v.SetDefault("environment", "development")
	v.SetDefault("log_level", "info")

	// Recipe API
	v.SetDefault("recipe_api.http_address", ":8081")
	v.SetDefault("recipe_api.grpc_address", ":9091")
	v.SetDefault("recipe_api.timeout", "30s")

	// MealPlanner API
	v.SetDefault("mealplanner_api.http_address", ":8082")
	v.SetDefault("mealplanner_api.grpc_address", ":9092")
	v.SetDefault("mealplanner_api.timeout", "30s")

	// BFF
	v.SetDefault("bff.http_address", ":8080")
	v.SetDefault("bff.timeout", "30s")
	v.SetDefault("bff.recipe_api_address", "localhost:9091")
	v.SetDefault("bff.mealplan_api_address", "localhost:9092")
	v.SetDefault("bff.cors_allowed_origins", []string{"*"})

	// Database
	v.SetDefault("database.recipe_db", "postgres://platepilot:platepilot@localhost:5432/recipedb?sslmode=disable")
	v.SetDefault("database.mealplanner_db", "postgres://platepilot:platepilot@localhost:5432/mealplannerdb?sslmode=disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_life", "5m")

	// RabbitMQ
	v.SetDefault("rabbitmq.url", "amqp://platepilot:platepilot@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "recipe-events")
	v.SetDefault("rabbitmq.exchange_type", "topic")

	// LLM (defaults for local Ollama development)
	v.SetDefault("llm.base_url", "http://localhost:11434/v1")
	v.SetDefault("llm.api_key", "ollama")
	v.SetDefault("llm.model", "llama3.1")
	v.SetDefault("llm.embed_model", "nomic-embed-text")
	v.SetDefault("llm.embed_dimensions", 768)
	v.SetDefault("llm.timeout", "60s")
	v.SetDefault("llm.max_tokens", 2000)
	v.SetDefault("llm.temperature", 0.7)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
