package crypto

import "golang.org/x/crypto/bcrypt"

func BcryptHash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compare compares the password with the hash
func BcryptCompare(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
