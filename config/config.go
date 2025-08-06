package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	User            string `env:"POSTGRES_USER"`
	Password        string `env:"POSTGRES_PASSWORD"`
	DBName          string `env:"POSTGRES_DB"`
	Port            int64  `env:"POSTGRES_PORT"`
	Host            string `env:"POSTGRES_HOST"`
	MaxIdleConns    int    `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"10"`
	MaxOpenConns    int    `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"100"`
	ConnMaxLifetime int64  `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"3600"`
}

type RedisConfig struct {
	Host string `env:"REDIS_HOST"`
	Port string `env:"REDIS_PORT"`
}

type BotConfig struct {
	Token           string `env:"TELEGRAM_BOT_TOKEN"`
	Admin           int64  `env:"TELEGRAM_ADMIN_USER_ID"`
	ChannelLogPanel int64  `env:"TELEGRAM_LOG_PANEL_VPN"`
}

type APIConfig struct {
	Host        string `env:"API_HOST"`
	Port        string `env:"API_PORT"`
	Debug       bool   `env:"APP_DEBUG"`
	Timezone    string `env:"APP_TIMEZONE"`
	FrontendURL string `env:"FRONTEND_URL"`
	JwtToken    string `env:"JWT_SECRET_KEY"`
}

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Bot      BotConfig
	API      APIConfig
}

var conf *Config

func GetConfig() *Config {

	return conf
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}

	// Parse each section into nested struct
	if err := env.Parse(&cfg.Postgres); err != nil {
		log.Fatalf("Failed to parse Postgres config: %v", err)
	}
	if err := env.Parse(&cfg.Redis); err != nil {
		log.Fatalf("Failed to parse Redis config: %v", err)
	}
	if err := env.Parse(&cfg.Bot); err != nil {
		log.Fatalf("Failed to parse Bot config: %v", err)
	}
	if err := env.Parse(&cfg.API); err != nil {
		log.Fatalf("Failed to parse API config: %v", err)
	}

	conf = cfg
	return conf
}
