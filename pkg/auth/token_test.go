package auth_test

import (
	"testing"
	"time"

	"github.com/HotPotatoC/roadmap_gen/pkg/auth"
	"github.com/HotPotatoC/roadmap_gen/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Token(t *testing.T) {
	t.Run("CreateToken", func(t *testing.T) {
		id := 1
		email := "test@example.com"

		token, err := auth.CreateToken(id, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("VerifyToken", func(t *testing.T) {
		id := 1
		email := "test@example.com"

		token, err := auth.CreateToken(id, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		payload, err := auth.VerifyToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, payload)
		assert.Equal(t, id, payload.ID)
		assert.Equal(t, email, payload.Email)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		id := 1
		email := "test@example.com"

		// Create a token with a short expiration time
		config.SetJWTSecretExpiresIn(time.Second * 1)
		token, err := auth.CreateToken(id, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Wait for the token to expire
		time.Sleep(time.Second * 2)

		payload, err := auth.VerifyToken(token)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})
}
