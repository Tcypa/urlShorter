package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

type Config struct {
	PostgresURL string `yaml:"postgresUrl"`
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("failed to decode config file: %w", err)
	}
}

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config not initialized. Call LoadConfig() first.")
	}
	return config
}
