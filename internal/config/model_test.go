package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
						Host:   "example.com",
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						Port:   22,
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout: 30,
				},
			},
		},
		{
			name: "preserves custom port",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						Port:   2222,
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						Port:   2222,
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout: 30,
				},
			},
		},
		{
			name: "sets default timeout",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						Port:   22,
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout: 30,
				},
			},
		},
		{
			name: "preserves custom timeout",
			config: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout: 60,
				},
			},
			expected: &Config{
				Profiles: map[string]*Profile{
					"test": {
						Host:   "example.com",
						Port:   22,
						User:   "user",
						Prompt: "$ ",
						Auth:   &Auth{Type: "password", Prompt: true},
					},
				},
				Options: &Options{
					Timeout: 60,
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
