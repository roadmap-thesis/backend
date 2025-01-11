package bcrypt

import baseBcrypt "golang.org/x/crypto/bcrypt"

func Hash(plain string) (string, error) {
	hash, err := baseBcrypt.GenerateFromPassword([]byte(plain), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compare compares the password with the hash
func Compare(hash, plain string) bool {
	err := baseBcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
