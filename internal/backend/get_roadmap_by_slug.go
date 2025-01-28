package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapOutput, error) {
	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}

	personalizationOpts, err := b.repository.PersonalizationOptions.GetByRoadmapID(ctx, roadmap.ID)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}

	account, err := b.repository.Account.GetByID(ctx, roadmap.AccountID)
	if err != nil {
		return io.GetRoadmapOutput{}, err
	}
	roadmap.SetCreator(account)

	output := new(io.GetRoadmapOutput)
	output.Make(roadmap, personalizationOpts)

	return *output, nil
}
