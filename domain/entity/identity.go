package entity

import (
	"time"

	"github.com/HotPotatoC/roadmap_gen/domain/object"
)

const (
	IdentityTable = "identities"
)

type Identity struct {
	ID       int
	Name     string
	Email    string
	Password object.Password

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewIdentity(name, email, plainPassword string) (*Identity, error) {
	password := object.Password(plainPassword)

	if err := password.Validate(plainPassword); err != nil {
		return nil, err
	}

	hash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	identity := &Identity{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	return identity, nil
}

func (e *Identity) CheckPassword(password string) bool {
	return e.Password.Compare(password)
}

func (e *Identity) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
