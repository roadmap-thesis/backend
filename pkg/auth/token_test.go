package auth_test

import (
	"testing"
	"time"

	"github.com/roadmap-thesis/backend/pkg/auth"
	"github.com/roadmap-thesis/backend/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Token(t *testing.T) {
	t.Parallel()
	t.Run("CreateToken", func(t *testing.T) {
		t.Parallel()
		id := 1

		token, err := auth.CreateToken(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("VerifyToken", func(t *testing.T) {
		t.Parallel()
		id := 1

		token, err := auth.CreateToken(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		payload, err := auth.VerifyToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, payload)
		assert.Equal(t, id, payload.ID)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		t.Parallel()
		id := 1

		// Create a token with a short expiration time
		config.SetJWTSecretExpiresIn(time.Second * 1)
		token, err := auth.CreateToken(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Wait for the token to expire
		time.Sleep(time.Second * 2)

		payload, err := auth.VerifyToken(token)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})
}
