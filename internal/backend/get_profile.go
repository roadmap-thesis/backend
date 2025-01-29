package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
	"go.opentelemetry.io/otel/attribute"
)

func (b *backend) GetProfile(ctx context.Context) (io.GetProfileOutput, error) {
	ctx, span := tracer.Start(ctx, "(*backend.GetProfile)")
	defer span.End()

	auth := auth.FromContext(ctx)

	span.SetAttributes(attribute.Int("account_id", auth.ID))

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
