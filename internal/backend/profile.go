package backend

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/pkg/auth"
	"github.com/HotPotatoC/roadmap_gen/pkg/commonerrors"
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
		return ProfileOutput{}, commonerrors.ResourceNotFound("Account")
	}

	return ProfileOutput{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}
