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
		{
			name: "valid multiple routes config",
			file: "../../test/fixtures/valid/multiple-routes.yml",
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

func TestValidate_EmptyRoutes(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/empty-routes.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "routes must have at least one route")
}

func TestValidate_InvalidRouteName(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/invalid-route-name.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contains invalid characters")
}

func TestValidate_MissingRoute(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/empty-route.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "route 'test-route' must have at least one step")
}

func TestValidate_MissingPasswordPrompt(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/invalid/missing-password-prompt.yml")
	require.NoError(t, err)

	err = Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password_prompt is required for password auth in route step 2")
}

func TestValidate_PasswordPromptInFirstStep(t *testing.T) {
	// password_prompt in 1st step should be allowed but ignored
	cfg := &Config{
		Version: "1.0",
		Profiles: map[string]*Profile{
			"bastion": {
				Host:         "bastion.example.com",
				User:         "user1",
				PromptMarker: "$ ",
				Auth: &Auth{
					Type:           "password",
					PasswordFile:   "passwords.dat",
					PasswordPrompt: "password:", // Allowed in 1st step (ignored)
				},
			},
		},
		Routes: map[string][]*RouteStep{
			"test-route": {
				{Profile: "bastion"},
			},
		},
	}

	err := Validate(cfg)
	assert.NoError(t, err, "password_prompt should be allowed in 1st step profile")
}

func TestValidate_PasswordPromptOnKeyfileAuth(t *testing.T) {
	// password_prompt on keyfile auth should be an error
	cfg := &Config{
		Version: "1.0",
		Profiles: map[string]*Profile{
			"server": {
				Host:         "server.example.com",
				User:         "user1",
				PromptMarker: "$ ",
				Auth: &Auth{
					Type:           "keyfile",
					Path:           "~/.ssh/id_rsa",
					PasswordPrompt: "password:", // Should not be set for keyfile
				},
			},
		},
		Routes: map[string][]*RouteStep{
			"test-route": {
				{Profile: "server"},
			},
		},
	}

	err := Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password_prompt should not be set for keyfile auth")
}

func TestValidate_PasswordPromptWithSingleQuote(t *testing.T) {
	// password_prompt with single quote should be an error (TTL injection prevention)
	cfg := &Config{
		Version: "1.0",
		Profiles: map[string]*Profile{
			"bastion": {
				Host:         "bastion.example.com",
				User:         "user1",
				PromptMarker: "$ ",
				Auth:         &Auth{Type: "password", PasswordFile: "passwords.dat"},
			},
			"target": {
				Host:         "target.example.com",
				User:         "user2",
				PromptMarker: "$ ",
				Auth: &Auth{
					Type:           "password",
					PasswordFile:   "passwords.dat",
					PasswordPrompt: "password':", // Single quote should be rejected
				},
			},
		},
		Routes: map[string][]*RouteStep{
			"test-route": {
				{Profile: "bastion"},
				{Profile: "target"},
			},
		},
	}

	err := Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password_prompt cannot contain single quotes")
}

func TestValidate_MultiHopMixedAuth(t *testing.T) {
	// Multi-hop route with mixed auth types: password -> keyfile -> password
	cfg := &Config{
		Version: "1.0",
		Profiles: map[string]*Profile{
			"bastion": {
				Host:         "bastion.example.com",
				User:         "user1",
				PromptMarker: "$ ",
				Auth:         &Auth{Type: "password", PasswordFile: "passwords.dat"},
			},
			"jump": {
				Host:         "jump.internal",
				User:         "user2",
				PromptMarker: "$ ",
				Auth:         &Auth{Type: "keyfile", Path: "~/.ssh/id_rsa"},
			},
			"target": {
				Host:         "target.internal",
				User:         "user3",
				PromptMarker: "$ ",
				Auth: &Auth{
					Type:           "password",
					PasswordFile:   "passwords.dat",
					PasswordPrompt: "password:",
				},
			},
		},
		Routes: map[string][]*RouteStep{
			"test-route": {
				{Profile: "bastion"}, // 1st step: password (no password_prompt needed)
				{Profile: "jump"},    // 2nd step: keyfile (no password_prompt needed)
				{Profile: "target"},  // 3rd step: password (password_prompt required)
			},
		},
	}

	err := Validate(cfg)
	assert.NoError(t, err, "multi-hop with mixed auth types should be valid")
}

func TestValidateAuth_Password(t *testing.T) {
	tests := []struct {
		name    string
		auth    *Auth
		wantErr bool
	}{
		{
			name:    "password with value",
			auth:    &Auth{Type: "password", Value: "secret"},
			wantErr: false,
		},
		{
			name:    "password with password_file",
			auth:    &Auth{Type: "password", PasswordFile: "passwords.dat"},
			wantErr: false,
		},
		{
			name:    "password without any source (gets default password_file)",
			auth:    &Auth{Type: "password"},
			wantErr: false, // デフォルト値が設定されるため成功
		},
		{
			name:    "password with both value and password_file",
			auth:    &Auth{Type: "password", Value: "secret", PasswordFile: "passwords.dat"},
			wantErr: true, // 相互排他エラー
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

func TestIsValidFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid names
		{"simple name", "simple", true},
		{"with hyphen", "multi-hop", true},
		{"with underscore", "route_1", true},
		{"mixed", "prod-db_v1", true},
		{"alphanumeric", "route123", true},

		// Invalid names
		{"empty string", "", false},
		{"with slash", "route/slash", false},
		{"with backslash", "route\\slash", false},
		{"path traversal", "../parent", false},
		{"absolute path", "/etc/passwd", false},
		{"with space", "route with spaces", false},
		{"with at", "route@special", false},
		{"with dot", "route.name", false},
		{"with colon", "route:name", false},
		{"with semicolon", "route;name", false},
		{"with pipe", "route|name", false},
		{"with ampersand", "route&name", false},
		{"with asterisk", "route*name", false},
		{"with question", "route?name", false},
		{"with quote", "route'name", false},
		{"with double quote", "route\"name", false},
		{"with angle bracket", "route<name", false},
		{"with angle bracket", "route>name", false},
		{"multibyte", "ルート", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidFileName(tt.input)
			assert.Equal(t, tt.expected, result, "isValidFileName(%q) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}
