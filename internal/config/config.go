package config

import "time"

// Config is the global config for the application
type Config struct {
	AppName string
	AppEnv  string
	Port    string

	DatabaseURL string

	JWTSecretKey       string
	JWTSecretExpiresIn time.Duration

	OpenAiAPIKey string
	OpenAiModel  string
}

var config *Config

// Init initializes the config package
func Init() {
	config = &Config{
		AppName: LookupEnv("APP_NAME", "roadmap_gen_api"),
		AppEnv:  LookupEnv("APP_ENV", "dev"),
		Port:    LookupEnv("PORT", "5000"),

		DatabaseURL: LookupEnv("DATABASE_URL", ""),

		JWTSecretKey:       LookupEnv("JWT_SECRET_KEY", "secret"),
		JWTSecretExpiresIn: LookupEnv("JWT_SECRET_EXPIRES_IN", time.Hour*24),

		OpenAiAPIKey: LookupEnv("OPENAI_API_KEY", ""),
		OpenAiModel:  LookupEnv("OPENAI_MODEL", "gpt-4o-mini"),
	}
}

// GetConfig returns the global config
func GetConfig() *Config { return config }

func AppName() string { return config.AppName }
func AppEnv() string  { return config.AppEnv }
func Port() string    { return config.Port }

func DatabaseURL() string { return config.DatabaseURL }

func OpenAiAPIKey() string { return config.OpenAiAPIKey }
func OpenAiModel() string  { return config.OpenAiModel }

func JWTSecretKey() string              { return config.JWTSecretKey }
func JWTSecretExpiresIn() time.Duration { return config.JWTSecretExpiresIn }
