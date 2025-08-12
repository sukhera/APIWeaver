package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/sukhera/APIWeaver/pkg/errors"
)

// Config represents the application configuration
type Config struct {
	// Parser settings
	StrictMode           bool          `mapstructure:"strict_mode" json:"strict_mode"`
	EnableRecovery       bool          `mapstructure:"enable_recovery" json:"enable_recovery"`
	MaxRecoveryAttempts  int           `mapstructure:"max_recovery_attempts" json:"max_recovery_attempts"`
	ParserTimeout        time.Duration `mapstructure:"parser_timeout" json:"parser_timeout"`
	InitialSliceCapacity int           `mapstructure:"initial_slice_capacity" json:"initial_slice_capacity"`

	// Validation settings
	ValidationLevel string   `mapstructure:"validation_level" json:"validation_level"`
	AllowedMethods  []string `mapstructure:"allowed_methods" json:"allowed_methods"`
	RequireExamples bool     `mapstructure:"require_examples" json:"require_examples"`
	MaxNestingDepth int      `mapstructure:"max_nesting_depth" json:"max_nesting_depth"`

	// Logging and monitoring
	Verbose         bool `mapstructure:"verbose" json:"verbose"`
	EnableMetrics   bool `mapstructure:"enable_metrics" json:"enable_metrics"`
	EnableProfiling bool `mapstructure:"enable_profiling" json:"enable_profiling"`

	// Output settings
	OutputFormat string `mapstructure:"output_format" json:"output_format"`
	PrettyPrint  bool   `mapstructure:"pretty_print" json:"pretty_print"`
}

// NewViperConfig creates a new Viper instance with default configuration
func NewViperConfig() *viper.Viper {
	v := viper.New()

	// Set default values
	v.SetDefault("strict_mode", false)
	v.SetDefault("enable_recovery", true)
	v.SetDefault("max_recovery_attempts", 3)
	v.SetDefault("parser_timeout", 30*time.Second)
	v.SetDefault("initial_slice_capacity", 100)
	v.SetDefault("validation_level", "basic")
	v.SetDefault("allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	v.SetDefault("require_examples", false)
	v.SetDefault("max_nesting_depth", 10)
	v.SetDefault("verbose", false)
	v.SetDefault("enable_metrics", false)
	v.SetDefault("enable_profiling", false)
	v.SetDefault("output_format", "json")
	v.SetDefault("pretty_print", true)

	// Configure Viper
	v.SetConfigName("apiweaver")        // name of config file (without extension)
	v.SetConfigType("yaml")             // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath(".")                // look for config in the working directory
	v.AddConfigPath("$HOME/.apiweaver") // call multiple times to add many search paths
	v.AddConfigPath("/etc/apiweaver/")  // path to look for the config file in

	// Environment variables
	v.SetEnvPrefix("APIWEAVER")
	v.AutomaticEnv() // read in environment variables that match

	// Bind environment variables
	_ = v.BindEnv("strict_mode", "APIWEAVER_STRICT_MODE")
	_ = v.BindEnv("verbose", "APIWEAVER_VERBOSE")
	_ = v.BindEnv("output_format", "APIWEAVER_OUTPUT_FORMAT")
	_ = v.BindEnv("parser_timeout", "APIWEAVER_PARSER_TIMEOUT")

	return v
}

// FromViper creates a Config from a Viper instance
func FromViper(v *viper.Viper) *Config {
	cfg := &Config{}

	// Unmarshal the configuration
	if err := v.Unmarshal(cfg); err != nil {
		// Return default config if unmarshaling fails
		return Default()
	}

	return cfg
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		StrictMode:           false,
		EnableRecovery:       true,
		MaxRecoveryAttempts:  3,
		ParserTimeout:        30 * time.Second,
		InitialSliceCapacity: 100,
		ValidationLevel:      "basic",
		AllowedMethods:       []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		RequireExamples:      false,
		MaxNestingDepth:      10,
		Verbose:              false,
		EnableMetrics:        false,
		EnableProfiling:      false,
		OutputFormat:         "json",
		PrettyPrint:          true,
	}
}

// LoadFromFile loads configuration from a YAML file using Viper
func LoadFromFile(filename string) (*Config, error) {
	v := NewViperConfig()
	v.SetConfigFile(filename)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := FromViper(v)
	return cfg, nil
}

// SaveToFile saves the configuration to a YAML file
func (c *Config) SaveToFile(filename string) error {
	v := viper.New()

	// Set values from config
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

	// Set config file
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")

	// Write config file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.MaxRecoveryAttempts < 0 {
		return errors.NewConfigError("max_recovery_attempts must be non-negative")
	}

	if c.ParserTimeout <= 0 {
		return errors.NewConfigError("parser_timeout must be positive")
	}

	if c.InitialSliceCapacity <= 0 {
		return errors.NewConfigError("initial_slice_capacity must be positive")
	}

	if c.MaxNestingDepth < 1 || c.MaxNestingDepth > 100 {
		return errors.NewConfigError("max_nesting_depth must be between 1 and 100")
	}

	validLevels := []string{"basic", "strict", "pedantic"}
	valid := false
	for _, level := range validLevels {
		if c.ValidationLevel == level {
			valid = true
			break
		}
	}
	if !valid {
		return errors.NewConfigError(fmt.Sprintf("validation_level must be one of: %v", validLevels))
	}

	if len(c.AllowedMethods) == 0 {
		return errors.NewConfigError("at least one HTTP method must be allowed")
	}

	validFormats := []string{"json", "yaml", "text"}
	valid = false
	for _, format := range validFormats {
		if c.OutputFormat == format {
			valid = true
			break
		}
	}
	if !valid {
		return errors.NewConfigError(fmt.Sprintf("output_format must be one of: %v", validFormats))
	}

	return nil
}

// ToViper converts the config to a Viper instance
func (c *Config) ToViper() *viper.Viper {
	v := viper.New()

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

	return v
}

// ToParserOptions converts the config to parser options
func (c *Config) ToParserOptions() []interface{} {
	return []interface{}{
		// These will be converted to actual parser options
		// when the parser package is properly integrated
	}
}
