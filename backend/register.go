package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/domain/entity"
)

type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (b *Backend) Register(ctx context.Context, input RegisterInput) error {
	identity, err := entity.NewIdentity(input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}

	_, err = b.repository.IdentityCreate(ctx, identity)
	if err != nil {
		return err
	}

	return nil
}
