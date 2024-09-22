package config

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)


type ServerConfig struct {
	Port     int
}

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
	Server ServerConfig
	Database DBConfig
}

// LoadConfig loads the configuration from config.yaml and environment variables
func LoadConfig() *Config {
    // Determine active profile (dev, prod, etc.) from environment variables or default to dev
    profile := os.Getenv("APP_PROFILE")
    if profile == "" {
        profile = "dev" // Default to dev profile
    }
    log.Printf("Loaded configuration profile: %s", profile)

    v := viper.New()

    // Load base config file
    v.SetConfigName("application") // Name of config file (without extension)
    v.SetConfigType("yaml")
    v.AddConfigPath("./config")    // Look for the config file in the "config" directory

    // Read the base config file (application.yaml)
    if err := v.ReadInConfig(); err != nil {
        log.Fatalf("Error loading base config file: %s", err)
    }

    // Load profile-specific config file (e.g., application-dev.yaml)
    v.SetConfigName("application-" + profile)
    if err := v.MergeInConfig(); err != nil {
        log.Fatalf("Error loading %s config file: %v", profile, err)
    }

    // Automatically map environment variables (e.g., SERVER_PORT → server.port)
    v.AutomaticEnv()

    // Replace "." with "_" for environment variable lookup
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Unmarshal the config into a struct
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        log.Fatalf("Error parsing config file: %s", err)
    }

    // Override struct values with any environment variables
    bindEnvVars(v, &cfg)

    return &cfg
}

// bindEnvVars uses reflection to bind environment variables to the config struct fields
func bindEnvVars(v *viper.Viper, config interface{}) {
    val := reflect.ValueOf(config).Elem() // Get the underlying value of the config struct

    bindFields(v, val, "")
}

// bindFields is a recursive function to traverse the struct and bind environment variables
func bindFields(v *viper.Viper, val reflect.Value, parentKey string) {
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        fieldName := fieldType.Name

        // Convert the field name into a configuration key (lowercase and dot notation)
        configKey := strings.ToLower(parentKey + fieldName)

        // If the field is a struct, recurse
        if field.Kind() == reflect.Struct {
            bindFields(v, field, configKey+".")
        } else {
            // Generate the ENV key (e.g., SERVER_PORT for configKey "server.port")
            envKey := strings.ToUpper(strings.Replace(configKey, ".", "_", -1))

            // Check if the environment variable exists
            if envVal := os.Getenv(envKey); envVal != "" {
                // Override the config struct field with the environment variable value
                v.Set(configKey, envVal)
            }
        }
    }
}
