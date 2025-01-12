package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
	"github.com/HotPotatoC/roadmap_gen/internal/domain/entity"
	"github.com/HotPotatoC/roadmap_gen/internal/jwt"
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

func (b *Backend) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
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

func (b *Backend) registerEmail(ctx context.Context, input RegisterInput) (*registerEmailOutput, error) {
	existingIdentity, err := b.repository.IdentityGetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	// sign in if identity already exists
	if existingIdentity != nil {
		matched := existingIdentity.CheckPassword(input.Password)

		if !matched {
			return nil, commonerrors.InvalidCredentials()
		}

		return &registerEmailOutput{id: existingIdentity.ID, email: existingIdentity.Email}, nil
	}

	identity, err := entity.NewIdentity(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	createdIdentity, err := b.repository.IdentityCreate(ctx, identity)
	if err != nil {
		return nil, err
	}

	return &registerEmailOutput{id: createdIdentity.ID, email: createdIdentity.Email, created: true}, nil
}
