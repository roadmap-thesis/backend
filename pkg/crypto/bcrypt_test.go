package crypto_test

import (
	"testing"

	"github.com/HotPotatoC/roadmap_gen/pkg/crypto"

	"github.com/stretchr/testify/assert"
)

func TestCrypto_BcryptHash(t *testing.T) {
	plainPassword := "mysecretpassword"
	hash, err := crypto.BcryptHash(plainPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCrypto_BcryptCompare(t *testing.T) {
	plainPassword := "mysecretpassword"
	hash, err := crypto.BcryptHash(plainPassword)
	assert.NoError(t, err)

	isValid := crypto.BcryptCompare(hash, plainPassword)
	assert.True(t, isValid)

	isInvalid := crypto.BcryptCompare(hash, "wrongpassword")
	assert.False(t, isInvalid)
}
