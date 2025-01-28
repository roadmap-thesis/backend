package domain

import (
	"time"

	"github.com/roadmap-thesis/backend/internal/domain/object"
)

const (
	AccountTable = "accounts"
)

type Account struct {
	ID       int
	Email    string
	Password object.Password

	Profile  *Profile
	Roadmaps []*Roadmap

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(email, plainPassword string) (*Account, error) {
	password := object.Password(plainPassword)

	if err := password.Validate(plainPassword); err != nil {
		return nil, err
	}

	hash, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account, nil
}

func (e *Account) IsZero() bool {
	return e.ID == 0 &&
		e.Email == "" &&
		e.Password == "" &&
		e.CreatedAt.IsZero() &&
		e.UpdatedAt.IsZero()
}

func (e *Account) CheckPassword(password string) bool {
	return e.Password.Compare(password)
}

func (e *Account) SetProfile(profile *Profile) {
	e.Profile = profile
}

func (e *Account) Update(name, email string) {
	e.Email = email
	e.UpdateChangelog()
}

func (e *Account) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
