package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/HotPotatoC/roadmap_gen/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Payload(t *testing.T) {
	t.Parallel()
	t.Run("NewPayload", func(t *testing.T) {
		t.Parallel()
		id := 1
		email := "test@example.com"
		expiresIn := time.Hour

		payload := auth.NewPayload(id, email, expiresIn)
		assert.NotNil(t, payload)
		assert.Equal(t, id, payload.ID)
		assert.Equal(t, email, payload.Email)
		assert.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
		assert.WithinDuration(t, time.Now().Add(expiresIn), payload.ExpiresAt, time.Second)
	})

	t.Run("NewFromClaims", func(t *testing.T) {
		t.Parallel()
		claims := jwt.MapClaims{
			"id":    float64(1),
			"email": "test@example.com",
			"iat":   float64(time.Now().Unix()),
			"exp":   float64(time.Now().Add(time.Hour).Unix()),
		}

		payload := auth.NewPayloadFromClaims(claims)
		assert.NotNil(t, payload)
		assert.Equal(t, int(claims["id"].(float64)), payload.ID)
		assert.Equal(t, claims["email"].(string), payload.Email)
		assert.WithinDuration(t, time.Unix(int64(claims["iat"].(float64)), 0), payload.IssuedAt, time.Second)
		assert.WithinDuration(t, time.Unix(int64(claims["exp"].(float64)), 0), payload.ExpiresAt, time.Second)
	})

	t.Run("Claims", func(t *testing.T) {
		t.Parallel()
		id := 1
		email := "test@example.com"
		expiresIn := time.Hour

		payload := auth.NewPayload(id, email, expiresIn)
		claims := payload.Claims().(jwt.MapClaims)

		assert.Equal(t, id, claims["id"])
		assert.Equal(t, email, claims["email"])
		assert.Equal(t, payload.IssuedAt.Unix(), claims["iat"])
		assert.Equal(t, payload.ExpiresAt.Unix(), claims["exp"])
	})

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		id := 1
		email := "test@example.com"
		expiresIn := time.Hour

		payload := auth.NewPayload(id, email, expiresIn)
		err := payload.Valid()
		assert.NoError(t, err)

		// Test expired payload
		expiredPayload := auth.NewPayload(id, email, -time.Hour)
		err = expiredPayload.Valid()
		assert.Error(t, err)
	})

	t.Run("FromContext", func(t *testing.T) {
		t.Parallel()
		id := 1
		email := "test@example.com"
		expiresIn := time.Hour

		payload := auth.NewPayload(id, email, expiresIn)
		ctx := context.WithValue(context.Background(), auth.AuthCtxKey, payload)

		extractedPayload := auth.FromContext(ctx)
		assert.NotNil(t, extractedPayload)
		assert.Equal(t, payload, extractedPayload)
	})
}
