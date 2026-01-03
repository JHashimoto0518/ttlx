package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_Success(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "valid simple config",
			file: "../../test/fixtures/valid/simple.yml",
		},
		{
			name: "valid full config",
			file: "../../test/fixtures/valid/full.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.file)
			require.NoError(t, err)

			err = Validate(cfg)
			assert.NoError(t, err)
		})
	}
}

func TestValidate_MissingVersion(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/missing-version.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "version field is required")
}

func TestValidate_MissingProfiles(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/missing-profiles.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one profile must be defined")
}

func TestValidate_InvalidProfileRef(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/invalid-profile-ref.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "profile 'nonexistent' not found")
}

func TestValidate_InvalidAuth(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/invalid-auth.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid auth type")
}

func TestValidate_MissingPromptMarker(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/missing-prompt-marker.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "prompt_marker is required")
}

func TestValidate_MissingRoute(t *testing.T) {
	cfg := &Config{
		Version: "1.0",
		Profiles: map[string]*Profile{
			"test": {
				Host:         "example.com",
				User:         "user",
				PromptMarker: "$ ",
				Auth:         &Auth{Type: "password", Prompt: true, PasswordPrompt: "password:"},
			},
		},
		Route: []*RouteStep{},
	}

	err := Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "route must have at least one step")
}

func TestValidate_MissingPasswordPrompt(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/missing-password-prompt.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password auth requires 'password_prompt'")
}

func TestValidateAuth_Password(t *testing.T) {
	tests := []struct {
		name    string
		auth    *Auth
		wantErr bool
	}{
		{
			name:    "password with value",
			auth:    &Auth{Type: "password", Value: "secret", PasswordPrompt: "password:"},
			wantErr: false,
		},
		{
			name:    "password with env",
			auth:    &Auth{Type: "password", Env: "PASSWORD_ENV", PasswordPrompt: "password:"},
			wantErr: false,
		},
		{
			name:    "password with prompt",
			auth:    &Auth{Type: "password", Prompt: true, PasswordPrompt: "password:"},
			wantErr: false,
		},
		{
			name:    "password without any source",
			auth:    &Auth{Type: "password", PasswordPrompt: "password:"},
			wantErr: true,
		},
		{
			name:    "password without password_prompt",
			auth:    &Auth{Type: "password", Prompt: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAuth(tt.auth)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAuth_Keyfile(t *testing.T) {
	tests := []struct {
		name    string
		auth    *Auth
		wantErr bool
	}{
		{
			name:    "keyfile with path",
			auth:    &Auth{Type: "keyfile", Path: "~/.ssh/id_rsa"},
			wantErr: false,
		},
		{
			name:    "keyfile without path",
			auth:    &Auth{Type: "keyfile"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAuth(tt.auth)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAuth_InvalidType(t *testing.T) {
	auth := &Auth{Type: "invalid"}
	err := validateAuth(auth)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid auth type")
}

func TestValidateAuth_Nil(t *testing.T) {
	err := validateAuth(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "auth is required")
}
