package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/HotPotatoC/roadmap_gen/internal/domain"
	"github.com/HotPotatoC/roadmap_gen/internal/jwt"
	"github.com/jackc/pgx/v5"
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

	token, err := jwt.GenerateJWT(map[string]any{
		"account_id":    result.id,
		"account_email": result.email,
	})
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
	var output *registerEmailOutput
	err := b.provider.Transaction.Execute(ctx, func(ctx context.Context, tx pgx.Tx) error {
		existingAccount, err := b.repository.Account.WithTx(tx).GetByEmail(ctx, input.Email)
		if err != nil {
			return err
		}

		// sign in if account already exists
		if existingAccount != nil {
			matched := existingAccount.CheckPassword(input.Password)

			if !matched {
				return commonerrors.InvalidCredentials()
			}

			output = &registerEmailOutput{id: existingAccount.ID, email: existingAccount.Email}
			return nil
		}

		account, err := domain.NewAccount(input.Name, input.Email, input.Password)
		if err != nil {
			return err
		}

		createdAccount, err := b.repository.Account.WithTx(tx).Create(ctx, account)
		if err != nil {
			return err
		}

		output = &registerEmailOutput{id: createdAccount.ID, email: createdAccount.Email, created: true}
		return nil
	})

	return output, err
}
