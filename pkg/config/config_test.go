package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/HotPotatoC/roadmap_gen/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setEnv(key, value string) {
	os.Setenv(key, value)
}

func unsetEnv(key string) {
	os.Unsetenv(key)
}

func TestConfig_Init(t *testing.T) {
	setEnv("APP_NAME", "test_app")
	setEnv("APP_ENV", "test")
	setEnv("PORT", "8080")
	setEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname")
	setEnv("DATABASE_MAX_CONNS", "10")
	setEnv("DATABASE_MIN_CONNS", "1")
	setEnv("DATABASE_MAX_CONN_LIFETIME", "2h")
	setEnv("DATABASE_MAX_CONN_IDLETIME", "1h")
	setEnv("DATABASE_HEALTH_CHECK_PERIOD", "30m")
	setEnv("DATABASE_DEFAULT_CONNECTION_TIMEOUT", "10s")
	setEnv("JWT_SECRET_KEY", "test_secret")
	setEnv("JWT_SECRET_EXPIRES_IN", "48h")
	setEnv("OPENAI_API_KEY", "test_api_key")
	setEnv("OPENAI_MODEL", "gpt-4")

	defer unsetEnv("APP_NAME")
	defer unsetEnv("APP_ENV")
	defer unsetEnv("PORT")
	defer unsetEnv("DATABASE_URL")
	defer unsetEnv("DATABASE_MAX_CONNS")
	defer unsetEnv("DATABASE_MIN_CONNS")
	defer unsetEnv("DATABASE_MAX_CONN_LIFETIME")
	defer unsetEnv("DATABASE_MAX_CONN_IDLETIME")
	defer unsetEnv("DATABASE_HEALTH_CHECK_PERIOD")
	defer unsetEnv("DATABASE_DEFAULT_CONNECTION_TIMEOUT")
	defer unsetEnv("JWT_SECRET_KEY")
	defer unsetEnv("JWT_SECRET_EXPIRES_IN")
	defer unsetEnv("OPENAI_API_KEY")
	defer unsetEnv("OPENAI_MODEL")

	config.Init()

	cfg := config.GetConfig()

	assert.Equal(t, "test_app", cfg.AppName)
	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/dbname", cfg.DatabaseURL)
	assert.Equal(t, int32(10), cfg.DatabaseMaxConns)
	assert.Equal(t, int32(1), cfg.DatabaseMinConns)
	assert.Equal(t, 2*time.Hour, cfg.DatabaseMaxConnLifetime)
	assert.Equal(t, 1*time.Hour, cfg.DatabaseMaxConnIdleTime)
	assert.Equal(t, 30*time.Minute, cfg.DatabaseHealthCheckPeriod)
	assert.Equal(t, 10*time.Second, cfg.DatabaseDefaultConnectionTimeout)
	assert.Equal(t, "test_secret", cfg.JWTSecretKey)
	assert.Equal(t, 48*time.Hour, cfg.JWTSecretExpiresIn)
	assert.Equal(t, "test_api_key", cfg.OpenAiAPIKey)
	assert.Equal(t, "gpt-4", cfg.OpenAiModel)
}

func TestConfig_SetOpenAiAPIKey(t *testing.T) {
	config.Init()
	config.SetOpenAiAPIKey("new_api_key")
	assert.Equal(t, "new_api_key", config.OpenAiAPIKey())
}

func TestConfig_SetOpenAiModel(t *testing.T) {
	config.Init()
	config.SetOpenAiModel("new_model")
	assert.Equal(t, "new_model", config.OpenAiModel())
}

func TestConfig_SetJWTSecretKey(t *testing.T) {
	config.Init()
	config.SetJWTSecretKey("new_secret_key")
	assert.Equal(t, "new_secret_key", config.JWTSecretKey())
}

func TestConfig_SetJWTSecretExpiresIn(t *testing.T) {
	config.Init()
	config.SetJWTSecretExpiresIn(72 * time.Hour)
	assert.Equal(t, 72*time.Hour, config.JWTSecretExpiresIn())
}
