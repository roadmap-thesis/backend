package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
)

func (b *backend) Profile(ctx context.Context) (io.ProfileOutput, error) {
	auth := auth.FromContext(ctx)

	account, err := b.repository.Account.GetByID(ctx, auth.ID)
	if err != nil {
		return io.ProfileOutput{}, err
	}

	return io.ProfileOutput{
		ID:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}, nil
}
