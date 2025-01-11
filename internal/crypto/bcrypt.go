package crypto

import baseBcrypt "golang.org/x/crypto/bcrypt"

type bcrypt struct {
	cost int
}

func NewBcrypt(cost int) Crypto {
	return bcrypt{cost: cost}
}

func (b bcrypt) Hash(plain string) (string, error) {
	hash, err := baseBcrypt.GenerateFromPassword([]byte(plain), b.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compare compares the password with the hash
func (b bcrypt) Compare(hash, plain string) bool {
	err := baseBcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
