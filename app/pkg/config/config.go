package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"sslmode"`
		Timezone string `yaml:"timezone"`
	} `yaml:"postgres"`
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"redis"`
	Server struct {
		Port           string `yaml:"port"`
		InternalSecret string `yaml:"internal_secret"`
		PublicSecret   string `yaml:"public_secret"`
	} `yaml:"server"`
}

func LoadConfig(filename string) *Config {
	raw, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	expanded := os.ExpandEnv(string(raw))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return &cfg
}
