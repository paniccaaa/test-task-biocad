package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	InputPath  string `yaml:"input_path"`
	OutputPath string `yaml:"output_path"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type ConfigDatabase struct {
	Port     string `env:"PORT"`
	Host     string `env:"HOST"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	User     string `env:"USER_NAME"`
	Password string `env:"PASSWORD"`
}

func MustLoad() (*Config, *ConfigDatabase) {
	configPath := "./config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error laoding .env file: %s", err)
	}

	var cfgDB ConfigDatabase
	if err := cleanenv.ReadEnv(&cfgDB); err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}

	return &cfg, &cfgDB
}
