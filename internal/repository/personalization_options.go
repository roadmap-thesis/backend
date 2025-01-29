package repository

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/pkg/database"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PersonalizationOptionsRepository struct {
	db database.Connection
}

func NewPersonalizationOptionsRepository(db database.Connection) *PersonalizationOptionsRepository {
	return &PersonalizationOptionsRepository{
		db: db,
	}
}

func (r *PersonalizationOptionsRepository) GetByRoadmapID(ctx context.Context, filter int) (domain.PersonalizationOptions, error) {
	ctx, span := tracer.Start(ctx, "(*PersonalizationOptionsRepository.GetByRoadmapID)", trace.WithAttributes(attribute.Int("roadmap_id", filter)))
	defer span.End()

	personalizationOpts, err := r.fetch(ctx, "roadmap_id", filter)
	if err != nil {
		return domain.PersonalizationOptions{}, err
	}

	if len(personalizationOpts) == 0 {
		return domain.PersonalizationOptions{}, domain.ErrPersonalizationOptionsNotFound
	}

	return personalizationOpts[0], nil
}

func (r *PersonalizationOptionsRepository) fetch(ctx context.Context, col string, args ...any) ([]domain.PersonalizationOptions, error) {
	ctx, span := tracer.Start(ctx, "(*PersonalizationOptionsRepository.fetch)", trace.WithAttributes(attribute.String("col", col)))
	defer span.End()

	query, args := psql.Select(
		sm.Columns(
			"id",
			"account_id",
			"roadmap_id",
			"daily_time_availability",
			"total_duration",
			"skill_level",
			"additional_info",
			"created_at",
			"updated_at",
		),
		sm.From(domain.PersonalizationOptionsTable),
		sm.Where(psql.Quote(col).EQ(psql.Arg(args...))),
	).MustBuild(ctx)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var personalizationOpts []domain.PersonalizationOptions
	for rows.Next() {
		var personalizationOpt domain.PersonalizationOptions
		err := rows.Scan(
			&personalizationOpt.ID,
			&personalizationOpt.AccountID,
			&personalizationOpt.RoadmapID,
			&personalizationOpt.DailyTimeAvailability,
			&personalizationOpt.TotalDuration,
			&personalizationOpt.SkillLevel,
			&personalizationOpt.AdditionalInfo,
			&personalizationOpt.CreatedAt,
			&personalizationOpt.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		personalizationOpts = append(personalizationOpts, personalizationOpt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return personalizationOpts, nil
}
