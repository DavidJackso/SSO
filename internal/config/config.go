package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Storage  string        `yaml:"storage" env-default:"s3"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPCConfig
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"4041"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
