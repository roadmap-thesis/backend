package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain/object"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
)

func (b *backend) ListUserRoadmaps(ctx context.Context) (io.ListUserRoadmapsOutput, error) {
	auth := auth.FromContext(ctx)

	roadmaps, err := b.repository.Roadmap.ListByAccountID(ctx, auth.ID)
	if err != nil {
		return io.ListUserRoadmapsOutput{}, err
	}

	output := io.ListUserRoadmapsOutput{
		TotalRoadmaps: len(roadmaps),
	}

	for _, roadmap := range roadmaps {
		outputRoadmap := io.ListUserRoadmapsOutputRoadmap{
			ID:                   roadmap.ID,
			Title:                roadmap.Title,
			Description:          roadmap.Description,
			Slug:                 roadmap.Slug,
			TotalTopics:          roadmap.TotalTopics(),
			CompletionPercentage: roadmap.CompletionPercentage(),
			PersonalizationOpts: io.ListUserRoadmapsOutputPersonalizationOptions{
				DailyTimeAvailability: object.NewIntervalFromDuration(roadmap.PersonalizationOptions.DailyTimeAvailability),
				TotalDuration:         object.NewIntervalFromDuration(roadmap.PersonalizationOptions.TotalDuration),
				SkillLevel:            roadmap.PersonalizationOptions.SkillLevel.String(),
				AdditionalInfo:        roadmap.PersonalizationOptions.AdditionalInfo,
			},
			CreatedAt: roadmap.CreatedAt,
			UpdatedAt: roadmap.UpdatedAt,
		}

		output.Roadmaps = append(output.Roadmaps, outputRoadmap)
	}

	return output, nil
}
