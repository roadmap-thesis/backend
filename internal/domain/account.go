package domain

import (
	"time"

	"github.com/HotPotatoC/roadmap_gen/internal/domain/object"
)

const (
	AccountTable = "accounts"
)

type Account struct {
	ID       int
	Name     string
	Email    string
	Password object.Password

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(name, email, plainPassword string) (*Account, error) {
	password := object.Password(plainPassword)

	if err := password.Validate(plainPassword); err != nil {
		return nil, err
	}

	hash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	return account, nil
}

func (e *Account) CheckPassword(password string) bool {
	return e.Password.Compare(password)
}

func (e *Account) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
