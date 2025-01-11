package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/domain/entity"
	"github.com/HotPotatoC/roadmap_gen/internal/jwt"
)

type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (b *Backend) Register(ctx context.Context, input RegisterInput) (string, error) {
	identity, err := entity.NewIdentity(input.Name, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	createdIdentity, err := b.repository.IdentityCreate(ctx, identity)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateJWT(map[string]any{
		"user_id":    createdIdentity.ID,
		"user_email": createdIdentity.Email,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
