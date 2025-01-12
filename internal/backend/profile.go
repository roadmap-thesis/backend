package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/auth"
	"github.com/HotPotatoC/roadmap_gen/internal/commonerrors"
)

type ProfileOutput struct {
	ID   int
	Name string
}

func (b *backend) Profile(ctx context.Context) (ProfileOutput, error) {
	identity := auth.FromContext(ctx)

	account, err := b.repository.Account.GetByID(ctx, identity.ID)
	if err != nil {
		return ProfileOutput{}, err
	}

	if account == nil {
		return ProfileOutput{}, commonerrors.NotFound("Account")
	}

	return ProfileOutput{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}
