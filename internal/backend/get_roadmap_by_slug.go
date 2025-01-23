package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/io"
)

func (b *backend) GetRoadmapBySlug(ctx context.Context, slug string) (io.GetRoadmapBySlugOutput, error) {
	roadmap, err := b.repository.Roadmap.GetBySlug(ctx, slug)
	if err != nil {
		return io.GetRoadmapBySlugOutput{}, err
	}

	account, err := b.repository.Account.GetByID(ctx, roadmap.AccountID)
	if err != nil {
		return io.GetRoadmapBySlugOutput{}, err
	}

	roadmap.SetCreator(account)

	output := new(io.GetRoadmapBySlugOutput)
	output.Make(roadmap)

	return *output, nil
}
