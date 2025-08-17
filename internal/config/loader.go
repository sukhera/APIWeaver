package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/sukhera/APIWeaver/internal/logger"
)

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port     int    `mapstructure:"port" json:"port"`
	Host     string `mapstructure:"host" json:"host"`
	DevMode  bool   `mapstructure:"dev_mode" json:"dev_mode"`
	CORS     CORSConfig
	Security SecurityConfig
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled" json:"enabled"`
	AllowedOrigins   []string `mapstructure:"allowed_origins" json:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" json:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials" json:"allow_credentials"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	RateLimiting RateLimitConfig `mapstructure:"rate_limiting" json:"rate_limiting"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool `mapstructure:"enabled" json:"enabled"`
	RequestsPerMinute int `mapstructure:"requests_per_minute" json:"requests_per_minute"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	MongoDB MongoDBConfig `mapstructure:"mongodb" json:"mongodb"`
	Cache   CacheConfig   `mapstructure:"cache" json:"cache"`
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	URI        string `mapstructure:"uri" json:"uri"`
	Database   string `mapstructure:"database" json:"database"`
	MaxPoolSize int   `mapstructure:"max_pool_size" json:"max_pool_size"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"` // seconds
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled     bool `mapstructure:"enabled" json:"enabled"`
	MaxSize     int  `mapstructure:"max_size" json:"max_size"`
	TTLSeconds  int  `mapstructure:"ttl_seconds" json:"ttl_seconds"`
}

// ExtendedConfig extends the base Config with additional fields
type ExtendedConfig struct {
	*Config
	Server  ServerConfig  `mapstructure:"server" json:"server"`
	Logger  logger.Config `mapstructure:"logger" json:"logger"`
	Storage StorageConfig `mapstructure:"storage" json:"storage"`
}

// Load loads configuration from file and environment variables
func Load(configFile string) (*ExtendedConfig, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Set config file
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Look for config file in standard locations
		v.SetConfigName("apiweaver")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.apiweaver")
		v.AddConfigPath("/etc/apiweaver")
	}

	// Environment variables
	v.SetEnvPrefix("APIWEAVER")
	v.AutomaticEnv()

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Initialize the base config first
	baseConfig := Default()
	
	// Unmarshal to struct  
	var cfg ExtendedConfig
	cfg.Config = baseConfig
	
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Base config defaults (from existing config.go)
	v.SetDefault("strict_mode", false)
	v.SetDefault("enable_recovery", true)
	v.SetDefault("max_recovery_attempts", 3)
	v.SetDefault("parser_timeout", "30s")
	v.SetDefault("initial_slice_capacity", 100)
	v.SetDefault("validation_level", "basic")
	v.SetDefault("allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	v.SetDefault("require_examples", false)
	v.SetDefault("max_nesting_depth", 10)
	v.SetDefault("verbose", false)
	v.SetDefault("enable_metrics", false)
	v.SetDefault("enable_profiling", false)
	v.SetDefault("output_format", "yaml")
	v.SetDefault("pretty_print", true)

	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.dev_mode", false)
	v.SetDefault("server.cors.enabled", true)
	v.SetDefault("server.cors.allowed_origins", []string{"*"})
	v.SetDefault("server.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("server.cors.allowed_headers", []string{"Content-Type", "Authorization"})
	v.SetDefault("server.cors.allow_credentials", false)
	v.SetDefault("server.security.rate_limiting.enabled", false)
	v.SetDefault("server.security.rate_limiting.requests_per_minute", 60)

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.output", "stdout")
	v.SetDefault("logger.add_source", false)

	// Storage defaults
	v.SetDefault("storage.mongodb.enabled", false)
	v.SetDefault("storage.mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("storage.mongodb.database", "apiweaver")
	v.SetDefault("storage.mongodb.max_pool_size", 10)
	v.SetDefault("storage.mongodb.timeout", 30)
	v.SetDefault("storage.cache.enabled", true)
	v.SetDefault("storage.cache.max_size", 1000)
	v.SetDefault("storage.cache.ttl_seconds", 3600)
}

// Save saves configuration to file
func (c *ExtendedConfig) Save(filename string) error {
	v := viper.New()

	// Set all values
	v.Set("strict_mode", c.StrictMode)
	v.Set("enable_recovery", c.EnableRecovery)
	v.Set("max_recovery_attempts", c.MaxRecoveryAttempts)
	v.Set("parser_timeout", c.ParserTimeout)
	v.Set("initial_slice_capacity", c.InitialSliceCapacity)
	v.Set("validation_level", c.ValidationLevel)
	v.Set("allowed_methods", c.AllowedMethods)
	v.Set("require_examples", c.RequireExamples)
	v.Set("max_nesting_depth", c.MaxNestingDepth)
	v.Set("verbose", c.Verbose)
	v.Set("enable_metrics", c.EnableMetrics)
	v.Set("enable_profiling", c.EnableProfiling)
	v.Set("output_format", c.OutputFormat)
	v.Set("pretty_print", c.PrettyPrint)

	// Server config
	v.Set("server", c.Server)
	v.Set("logger", c.Logger)
	v.Set("storage", c.Storage)

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}