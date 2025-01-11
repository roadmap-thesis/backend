package object

import (
	"errors"
	"unicode"

	"github.com/HotPotatoC/roadmap_gen/internal/crypto"
)

var (
	// ErrPasswordInvalidCharacters is returned when the password contains invalid characters
	ErrPasswordInvalidCharacters = errors.New("err password invalid characters")
)

// Password is your average user's secret password
type Password struct {
	hash   string
	crypto crypto.Crypto
}

func NewPassword() *Password {
	return &Password{
		crypto: crypto.NewBcrypt(10),
	}
}

func NewPasswordFrom(hash string) *Password {
	return &Password{hash: hash}
}

// Validate validates a plain password
func (p *Password) Validate(plain string) error {
	if !p.validateCharacters(plain) {
		return ErrPasswordInvalidCharacters
	}

	return nil
}

func (p *Password) validateCharacters(plain string) bool {
	for _, char := range plain {
		if char > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// GenerateHash generates a hash for the password
func (p *Password) GenerateHash(plain string) error {
	hash, err := p.crypto.Hash(plain)
	if err != nil {
		return err
	}

	p.hash = hash

	return nil
}

// Compare compares the password with the hash
func (p *Password) Compare(plain string) bool {
	return p.crypto.Compare(p.hash, plain)
}
