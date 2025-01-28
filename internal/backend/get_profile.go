package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
)

func (b *backend) GetProfile(ctx context.Context) (io.GetProfileOutput, error) {
	auth := auth.FromContext(ctx)

	account, err := b.repository.Account.GetByID(ctx, auth.ID)
	if err != nil {
		return io.GetProfileOutput{}, err
	}

	return io.GetProfileOutput{
		ID:       account.ID,
		Email:    account.Email,
		Name:     account.Profile.Name,
		Avatar:   account.Profile.Avatar,
		JoinedAt: account.CreatedAt,
	}, nil
}
