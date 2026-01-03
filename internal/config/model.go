package config

// Config represents the entire YAML configuration.
type Config struct {
	Version  string                       `yaml:"version"`
	Profiles map[string]*Profile          `yaml:"profiles"`
	Routes   map[string][]*RouteStep      `yaml:"routes"`
	Options  *Options                     `yaml:"options,omitempty"`
}

// Profile represents an SSH connection profile.
type Profile struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port,omitempty"`          // デフォルト: 22
	User         string `yaml:"user"`
	PromptMarker string `yaml:"prompt_marker"`           // プロンプトを識別する文字列（必須）例: "$ ", "# "
	Auth         *Auth  `yaml:"auth"`
}

// Auth represents authentication settings.
type Auth struct {
	Type           string `yaml:"type"`                      // "password" | "keyfile"
	Value          string `yaml:"value,omitempty"`           // パスワード直接記述
	Env            string `yaml:"env,omitempty"`             // 環境変数名
	Prompt         bool   `yaml:"prompt,omitempty"`          // 実行時入力
	PasswordPrompt string `yaml:"password_prompt,omitempty"` // パスワード入力待機文字列（2段目以降で必須）例: "password:"
	Path           string `yaml:"path,omitempty"`            // 秘密鍵ファイルパス
}

// RouteStep represents a step in the connection route.
type RouteStep struct {
	Profile  string   `yaml:"profile"`
	Commands []string `yaml:"commands,omitempty"`
}

// Options represents global options.
type Options struct {
	Timeout        int    `yaml:"timeout,omitempty"`
	Retry          int    `yaml:"retry,omitempty"`
	Log            bool   `yaml:"log,omitempty"`
	LogFile        string `yaml:"log_file,omitempty"`
	AutoDisconnect *bool  `yaml:"auto_disconnect,omitempty"` // 最終ステップ完了後に自動切断するか（デフォルト: false）
}

// SetDefaults sets default values for the config.
func (c *Config) SetDefaults() {
	for _, profile := range c.Profiles {
		if profile.Port == 0 {
			profile.Port = 22
		}
	}

	if c.Options == nil {
		c.Options = &Options{}
	}
	if c.Options.Timeout == 0 {
		c.Options.Timeout = 30
	}
	if c.Options.AutoDisconnect == nil {
		defaultAutoDisconnect := false
		c.Options.AutoDisconnect = &defaultAutoDisconnect
	}
}
