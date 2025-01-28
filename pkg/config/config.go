package config

import (
	"time"

	"github.com/cohesion-org/deepseek-go"
)

// Config is the global config for the application
type Config struct {
	AppName string
	AppEnv  string
	Port    string

	DatabaseURL                      string
	DatabaseMaxConns                 int32
	DatabaseMinConns                 int32
	DatabaseMaxConnLifetime          time.Duration
	DatabaseMaxConnIdleTime          time.Duration
	DatabaseHealthCheckPeriod        time.Duration
	DatabaseDefaultConnectionTimeout time.Duration

	JWTSecretKey       string
	JWTSecretExpiresIn time.Duration

	OpenAiAPIKey string
	OpenAiModel  string

	DeepSeekAPIKey string
	DeepSeekModel  string
}

var config *Config

// Init initializes the config package
func Init() {
	config = &Config{
		AppName: LookupEnv("APP_NAME", "roadmap_gen_api"),
		AppEnv:  LookupEnv("APP_ENV", "local"),
		Port:    LookupEnv("PORT", "5000"),

		DatabaseURL:                      LookupEnv("DATABASE_URL", ""),
		DatabaseMaxConns:                 LookupEnv("DATABASE_MAX_CONNS", int32(4)),
		DatabaseMinConns:                 LookupEnv("DATABASE_MIN_CONNS", int32(0)),
		DatabaseMaxConnLifetime:          LookupEnv("DATABASE_MAX_CONN_LIFETIME", time.Hour),
		DatabaseMaxConnIdleTime:          LookupEnv("DATABASE_MAX_CONN_IDLETIME", time.Minute*30),
		DatabaseHealthCheckPeriod:        LookupEnv("DATABASE_HEALTH_CHECK_PERIOD", time.Minute),
		DatabaseDefaultConnectionTimeout: LookupEnv("DATABASE_DEFAULT_CONNECTION_TIMEOUT", time.Second*5),

		JWTSecretKey:       LookupEnv("JWT_SECRET_KEY", "secret"),
		JWTSecretExpiresIn: LookupEnv("JWT_SECRET_EXPIRES_IN", time.Hour*24),

		OpenAiAPIKey: LookupEnv("OPENAI_API_KEY", ""),
		OpenAiModel:  LookupEnv("OPENAI_MODEL", "gpt-4o-mini"),

		DeepSeekAPIKey: LookupEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekModel:  LookupEnv("DEEPSEEK_MODEL", deepseek.DeepSeekChat),
	}
}

// GetConfig returns the global config
func GetConfig() *Config { return config }

func AppName() string { return config.AppName }
func AppEnv() string  { return config.AppEnv }
func Port() string    { return config.Port }

func DatabaseURL() string                             { return config.DatabaseURL }
func DatabaseMaxConns() int32                         { return config.DatabaseMaxConns }
func DatabaseMinConns() int32                         { return config.DatabaseMinConns }
func DatabaseMaxConnLifetime() time.Duration          { return config.DatabaseMaxConnLifetime }
func DatabaseMaxConnIdleTime() time.Duration          { return config.DatabaseMaxConnIdleTime }
func DatabaseHealthCheckPeriod() time.Duration        { return config.DatabaseHealthCheckPeriod }
func DatabaseDefaultConnectionTimeout() time.Duration { return config.DatabaseDefaultConnectionTimeout }

func JWTSecretKey() string              { return config.JWTSecretKey }
func JWTSecretExpiresIn() time.Duration { return config.JWTSecretExpiresIn }

func SetJWTSecretKey(key string)                   { config.JWTSecretKey = key }
func SetJWTSecretExpiresIn(duration time.Duration) { config.JWTSecretExpiresIn = duration }

func OpenAiAPIKey() string { return config.OpenAiAPIKey }
func OpenAiModel() string  { return config.OpenAiModel }

func SetOpenAiAPIKey(key string)  { config.OpenAiAPIKey = key }
func SetOpenAiModel(model string) { config.OpenAiModel = model }

func DeepSeekAPIKey() string { return config.DeepSeekAPIKey }
func DeepSeekModel() string  { return config.DeepSeekModel }

func SetDeepSeekAPIKey(key string)  { config.DeepSeekAPIKey = key }
func SetDeepSeekModel(model string) { config.DeepSeekModel = model }
