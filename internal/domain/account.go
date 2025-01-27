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
	Name     string
	Email    string
	Password object.Password

	Roadmaps []*Roadmap

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
		Name:      name,
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account, nil
}

func (e *Account) IsZero() bool {
	return e.ID == 0 &&
		e.Name == "" &&
		e.Email == "" &&
		e.Password == "" &&
		e.CreatedAt.IsZero() &&
		e.UpdatedAt.IsZero()
}

func (e *Account) CheckPassword(password string) bool {
	return e.Password.Compare(password)
}

func (e *Account) Update(name, email string) {
	e.Name = name
	e.Email = email
	e.UpdateChangelog()
}

func (e *Account) UpdateChangelog() {
	e.UpdatedAt = time.Now()
}
