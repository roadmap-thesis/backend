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

// Password can be plain/hashed
type Password string

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

// Hash generates a hash for the password
func (p Password) Hash(plain string) (Password, error) {
	hash, err := crypto.BcryptHash(plain)
	if err != nil {
		return "", err
	}

	return Password(hash), nil
}

// Compare compares the password with the hash
func (p Password) Compare(plain string) bool {
	return crypto.BcryptCompare(string(p), plain)
}

func (p Password) String() string {
	return string(p)
}
