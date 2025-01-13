package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/domain"
	"github.com/HotPotatoC/roadmap_gen/pkg/auth"
	"github.com/HotPotatoC/roadmap_gen/pkg/commonerrors"
)

type AuthInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthOutput struct {
	Created bool
	Token   string
}

func (b *backend) Auth(ctx context.Context, input AuthInput) (AuthOutput, error) {
	var output AuthOutput

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

func (b *backend) registerEmail(ctx context.Context, input AuthInput) (*registerEmailOutput, error) {
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
