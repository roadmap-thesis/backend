package crypto_test

import (
	"testing"

	"github.com/roadmap-thesis/backend/pkg/crypto"

	"github.com/stretchr/testify/assert"
)

func TestCrypto_BcryptHash(t *testing.T) {
	t.Parallel()
	plainPassword := "mysecretpassword"
	hash, err := crypto.BcryptHash(plainPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCrypto_BcryptCompare(t *testing.T) {
	t.Parallel()
	plainPassword := "mysecretpassword"
	hash, err := crypto.BcryptHash(plainPassword)
	assert.NoError(t, err)

	isValid := crypto.BcryptCompare(hash, plainPassword)
	assert.True(t, isValid)

	isInvalid := crypto.BcryptCompare(hash, "wrongpassword")
	assert.False(t, isInvalid)
}
