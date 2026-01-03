package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// boolPtr returns a pointer to a bool value.
func boolPtr(b bool) *bool {
	return &b
}

func TestConfig_SetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected *Config
	}{
		{
			name: "sets default port to 22",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						Port:         22,
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout:        30,
					AutoDisconnect: boolPtr(false),
				},
			},
		},
		{
			name: "preserves custom port",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						Port:         2222,
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						Port:         2222,
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout:        30,
					AutoDisconnect: boolPtr(false),
				},
			},
		},
		{
			name: "sets default timeout",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						Port:         22,
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout:        30,
					AutoDisconnect: boolPtr(false),
				},
			},
		},
		{
			name: "preserves custom timeout",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout:        60,
					AutoDisconnect: boolPtr(false),
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						Port:         22,
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout:        60,
					AutoDisconnect: boolPtr(false),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			assert.Equal(t, tt.expected, tt.config)
		})
	}
}

func TestConfig_SetDefaults_AutoDisconnect(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "auto_disconnect not specified (default to false)",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{},
			},
			expected: false,
		},
		{
			name: "auto_disconnect explicitly set to true",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					AutoDisconnect: boolPtr(true),
				},
			},
			expected: true,
		},
		{
			name: "auto_disconnect explicitly set to false",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:         "example.com",
						User:         "user",
						PromptMarker: "$ ",
						Auth:         &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					AutoDisconnect: boolPtr(false),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			assert.NotNil(t, tt.config.Options.AutoDisconnect)
			assert.Equal(t, tt.expected, *tt.config.Options.AutoDisconnect)
		})
	}
}
