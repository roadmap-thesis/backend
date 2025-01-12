package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/HotPotatoC/roadmap_gen/internal/domain/entity"
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
		"user_id":    result.id,
		"user_email": result.email,
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
		existingIdentity, err := b.repository.Identity.WithTx(tx).GetByEmail(ctx, input.Email)
		if err != nil {
			return err
		}

		// sign in if identity already exists
		if existingIdentity != nil {
			matched := existingIdentity.CheckPassword(input.Password)

			if !matched {
				return commonerrors.InvalidCredentials()
			}

			output = &registerEmailOutput{id: existingIdentity.ID, email: existingIdentity.Email}
			return nil
		}

		identity, err := entity.NewIdentity(input.Name, input.Email, input.Password)
		if err != nil {
			return err
		}

		createdIdentity, err := b.repository.Identity.WithTx(tx).Create(ctx, identity)
		if err != nil {
			return err
		}

		output = &registerEmailOutput{id: createdIdentity.ID, email: createdIdentity.Email, created: true}
		return nil
	})

	return output, err
}
