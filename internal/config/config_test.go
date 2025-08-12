package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Default(t *testing.T) {
	cfg := Default()

	assert.NotNil(t, cfg)
	assert.False(t, cfg.StrictMode)
	assert.True(t, cfg.EnableRecovery)
	assert.Equal(t, 3, cfg.MaxRecoveryAttempts)
	assert.Equal(t, 30*time.Second, cfg.ParserTimeout)
	assert.Equal(t, 100, cfg.InitialSliceCapacity)
	assert.Equal(t, "basic", cfg.ValidationLevel)
	assert.Len(t, cfg.AllowedMethods, 7)
	assert.False(t, cfg.RequireExamples)
	assert.Equal(t, 10, cfg.MaxNestingDepth)
	assert.False(t, cfg.Verbose)
	assert.False(t, cfg.EnableMetrics)
	assert.False(t, cfg.EnableProfiling)
	assert.Equal(t, "json", cfg.OutputFormat)
	assert.True(t, cfg.PrettyPrint)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "invalid max recovery attempts",
			config: func() *Config {
				cfg := Default()
				cfg.MaxRecoveryAttempts = -1
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "invalid parser timeout",
			config: func() *Config {
				cfg := Default()
				cfg.ParserTimeout = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "invalid initial slice capacity",
			config: func() *Config {
				cfg := Default()
				cfg.InitialSliceCapacity = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "invalid max nesting depth",
			config: func() *Config {
				cfg := Default()
				cfg.MaxNestingDepth = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "invalid validation level",
			config: func() *Config {
				cfg := Default()
				cfg.ValidationLevel = "invalid"
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "empty allowed methods",
			config: func() *Config {
				cfg := Default()
				cfg.AllowedMethods = []string{}
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "invalid output format",
			config: func() *Config {
				cfg := Default()
				cfg.OutputFormat = "invalid"
				return cfg
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewViperConfig(t *testing.T) {
	v := NewViperConfig()
	require.NotNil(t, v)

	// Test default values
	assert.Equal(t, false, v.GetBool("strict_mode"))
	assert.Equal(t, true, v.GetBool("enable_recovery"))
	assert.Equal(t, 3, v.GetInt("max_recovery_attempts"))
	assert.Equal(t, 30*time.Second, v.GetDuration("parser_timeout"))
	assert.Equal(t, 100, v.GetInt("initial_slice_capacity"))
	assert.Equal(t, "basic", v.GetString("validation_level"))
	assert.Equal(t, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, v.GetStringSlice("allowed_methods"))
	assert.Equal(t, false, v.GetBool("require_examples"))
	assert.Equal(t, 10, v.GetInt("max_nesting_depth"))
	assert.Equal(t, false, v.GetBool("verbose"))
	assert.Equal(t, false, v.GetBool("enable_metrics"))
	assert.Equal(t, false, v.GetBool("enable_profiling"))
	assert.Equal(t, "json", v.GetString("output_format"))
	assert.Equal(t, true, v.GetBool("pretty_print"))
}

func TestFromViper(t *testing.T) {
	v := viper.New()
	v.Set("strict_mode", true)
	v.Set("verbose", true)
	v.Set("output_format", "yaml")
	v.Set("parser_timeout", 60*time.Second)

	cfg := FromViper(v)
	require.NotNil(t, cfg)

	assert.True(t, cfg.StrictMode)
	assert.True(t, cfg.Verbose)
	assert.Equal(t, "yaml", cfg.OutputFormat)
	assert.Equal(t, 60*time.Second, cfg.ParserTimeout)
}

func TestConfig_LoadFromFile(t *testing.T) {
	// Create a temporary config file
	configContent := `
strict_mode: true
enable_recovery: false
max_recovery_attempts: 5
parser_timeout: 60s
initial_slice_capacity: 200
validation_level: strict
allowed_methods: ["GET", "POST"]
require_examples: true
max_nesting_depth: 15
verbose: true
enable_metrics: true
enable_profiling: true
output_format: yaml
pretty_print: false
`

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer func() {
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			t.Logf("Failed to remove temp file: %v", removeErr)
		}
	}()

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	// Test loading from file
	cfg, err := LoadFromFile(tmpFile.Name())
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify loaded values
	assert.True(t, cfg.StrictMode)
	assert.False(t, cfg.EnableRecovery)
	assert.Equal(t, 5, cfg.MaxRecoveryAttempts)
	assert.Equal(t, 60*time.Second, cfg.ParserTimeout)
	assert.Equal(t, 200, cfg.InitialSliceCapacity)
	assert.Equal(t, "strict", cfg.ValidationLevel)
	assert.Equal(t, []string{"GET", "POST"}, cfg.AllowedMethods)
	assert.True(t, cfg.RequireExamples)
	assert.Equal(t, 15, cfg.MaxNestingDepth)
	assert.True(t, cfg.Verbose)
	assert.True(t, cfg.EnableMetrics)
	assert.True(t, cfg.EnableProfiling)
	assert.Equal(t, "yaml", cfg.OutputFormat)
	assert.False(t, cfg.PrettyPrint)
}

func TestConfig_LoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("nonexistent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestConfig_LoadFromFile_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	configContent := `
strict_mode: true
invalid: yaml: content
`

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer func() {
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			t.Logf("Failed to remove temp file: %v", removeErr)
		}
	}()

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	// Test loading from file - should fail with invalid YAML
	_, err = LoadFromFile(tmpFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestConfig_SaveToFile(t *testing.T) {
	cfg := Default()
	cfg.StrictMode = true
	cfg.Verbose = true

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)
	defer func() {
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			t.Logf("Failed to remove temp file: %v", removeErr)
		}
	}()

	// Save config to file
	err = cfg.SaveToFile(tmpFile.Name())
	require.NoError(t, err)

	// Verify file was created and has content
	fileInfo, err := os.Stat(tmpFile.Name())
	require.NoError(t, err)
	assert.Greater(t, fileInfo.Size(), int64(0))

	// Load the saved config and verify it matches
	loadedCfg, err := LoadFromFile(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, cfg.StrictMode, loadedCfg.StrictMode)
	assert.Equal(t, cfg.Verbose, loadedCfg.Verbose)
}

func TestConfig_ToViper(t *testing.T) {
	cfg := Default()
	cfg.StrictMode = true
	cfg.Verbose = true
	cfg.OutputFormat = "yaml"

	v := cfg.ToViper()
	require.NotNil(t, v)

	assert.True(t, v.GetBool("strict_mode"))
	assert.True(t, v.GetBool("verbose"))
	assert.Equal(t, "yaml", v.GetString("output_format"))
}

func TestConfig_ValidationLevels(t *testing.T) {
	validLevels := []string{"basic", "strict", "pedantic"}

	for _, level := range validLevels {
		t.Run("valid_"+level, func(t *testing.T) {
			cfg := Default()
			cfg.ValidationLevel = level
			err := cfg.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestConfig_OutputFormats(t *testing.T) {
	validFormats := []string{"json", "yaml", "text"}

	for _, format := range validFormats {
		t.Run("valid_"+format, func(t *testing.T) {
			cfg := Default()
			cfg.OutputFormat = format
			err := cfg.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestConfig_AllowedMethods(t *testing.T) {
	tests := []struct {
		name    string
		methods []string
		wantErr bool
	}{
		{
			name:    "valid methods",
			methods: []string{"GET", "POST", "PUT"},
			wantErr: false,
		},
		{
			name:    "empty methods",
			methods: []string{},
			wantErr: true,
		},
		{
			name:    "nil methods",
			methods: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.AllowedMethods = tt.methods
			err := cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
