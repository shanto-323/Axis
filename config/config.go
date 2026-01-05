package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        Server               `koanf:"server" validate:"required"`
	Database      Database             `koanf:"database" validate:"required"`
	AiManage      AiManager            `koanf:"ai_manager" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type Server struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
	JwtKey             string   `koanf:"jwt_key" validate:"required"`
}

type Database struct {
	Type            string `koanf:"type" validate:"required,oneof=mock postgres"`
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type AiManager struct {
	Provider string `koanf:"provider" validate:"required"`
	ApiKey   string `koanf:"api_key" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	k := koanf.New(".")
	err := k.Load(env.Provider("", k.Delim(), func(s string) string {
		return strings.ToLower(s)
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	config := &Config{}
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling  %w", err)
	}

	if err := validator.New().Struct(config); err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main ")
	}

	if config.Observability == nil {
		config.Observability = DefaultObservabilityConfig()
	}

	config.Observability.ServiceName = "tasker"
	config.Observability.Environment = config.Primary.Env

	if err := config.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return config, nil
}

func (cfg *Config) IsProd() bool {
	return cfg.Primary.Env == "prod"
}
