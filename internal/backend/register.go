package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/auth"
	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/HotPotatoC/roadmap_gen/internal/domain"
)

type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterOutput struct {
	Created bool
	Token   string
}

func (b *backend) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	var output RegisterOutput

	result, err := b.registerEmail(ctx, input)
	if err != nil {
		return output, err
	}

	token, err := auth.CreateToken(result.id, result.email)
	if err != nil {
		return output, err
	}

	output.Created = result.created
	output.Token = token
	return output, nil
}

type registerEmailOutput struct {
	id      int
	email   string
	created bool
}

func (b *backend) registerEmail(ctx context.Context, input RegisterInput) (*registerEmailOutput, error) {
	existingAccount, err := b.repository.Account.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	// sign in if account already exists
	if existingAccount != nil {
		matched := existingAccount.CheckPassword(input.Password)

		if !matched {
			return nil, commonerrors.InvalidCredentials()
		}

		return &registerEmailOutput{id: existingAccount.ID, email: existingAccount.Email}, nil
	}

	account, err := domain.NewAccount(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	createdAccount, err := b.repository.Account.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &registerEmailOutput{id: createdAccount.ID, email: createdAccount.Email, created: true}, err
}
