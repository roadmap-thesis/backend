package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"github.com/roadmap-thesis/backend/pkg/auth"
)

func (b *backend) Auth(ctx context.Context, input io.AuthInput) (io.AuthOutput, error) {
	var output io.AuthOutput

	result, err := b.registerEmail(ctx, input)
	if err != nil {
		return output, err
	}

	token, err := auth.CreateToken(result.id)
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

func (b *backend) registerEmail(ctx context.Context, input io.AuthInput) (registerEmailOutput, error) {
	existingAccount, err := b.repository.Account.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrAccountNotFound {
		return registerEmailOutput{}, err
	}

	// sign in if account already exists
	if !existingAccount.IsZero() {
		matched := existingAccount.CheckPassword(input.Password)

		if !matched {
			return registerEmailOutput{}, apperrors.InvalidCredentials()
		}

		return registerEmailOutput{id: existingAccount.ID, email: existingAccount.Email}, nil
	}

	profile := domain.NewProfile(input.Name, input.Avatar)
	account, err := domain.NewAccount(input.Email, input.Password, profile)
	if err != nil {
		return registerEmailOutput{}, err
	}

	createdAccount, err := b.repository.Account.Save(ctx, account)
	if err != nil {
		return registerEmailOutput{}, err
	}

	return registerEmailOutput{id: createdAccount.ID, email: createdAccount.Email, created: true}, err
}
