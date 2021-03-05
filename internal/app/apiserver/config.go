package apiserver

// Config is a server configuration structure
type Config struct {
	BindAddr    string `yaml:"bind-addr"`
	Port        string `yaml:"port"`
	LogLevel    string `yaml:"log-level"`
	DatabaseURL string `yaml:"database-url"`
	SessionKey  string `yaml:"session-key"`
	CookieAge   int    `yaml:"cookie-age"`
}

// NewConfig is a constructor helper
func NewConfig() *Config {
	return &Config{
		BindAddr:  "localhost",
		Port:      "8081",
		LogLevel:  "debug",
		CookieAge: 600,
	}
}
