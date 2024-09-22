package config

import (
	"github.com/spf13/viper"
	"log"
)

// DBConfig holds the database configuration details
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Config stores all application configuration
type Config struct {
	Database DBConfig
}

// LoadConfig loads the configuration from config.yaml
func LoadConfig() *Config {
	var cfg Config

	viper.SetConfigName("config") // Name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // Look for the config file in the current directory

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	// Unmarshal the config file into the config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	return &cfg
}
