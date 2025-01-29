package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/pkg/database"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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

func (r *PersonalizationOptionsRepository) Save(ctx context.Context, input *domain.PersonalizationOptions) (domain.PersonalizationOptions, error) {
	query, args := psql.Insert(
		im.Into(domain.PersonalizationOptionsTable,
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
		im.Values(psql.Arg(
			input.ID,
			input.AccountID,
			input.RoadmapID,
			input.DailyTimeAvailability,
			input.TotalDuration,
			input.SkillLevel,
			input.AdditionalInfo,
			input.CreatedAt,
			input.UpdatedAt,
		)),
		im.Returning(
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
	).MustBuild(ctx)

	var personalizationOpts domain.PersonalizationOptions
	err := r.db.InTx(ctx, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, args...).Scan(
			&personalizationOpts.ID,
			&personalizationOpts.AccountID,
			&personalizationOpts.RoadmapID,
			&personalizationOpts.DailyTimeAvailability,
			&personalizationOpts.TotalDuration,
			&personalizationOpts.SkillLevel,
			&personalizationOpts.AdditionalInfo,
			&personalizationOpts.CreatedAt,
			&personalizationOpts.UpdatedAt,
		)
	})
	if err != nil {
		return domain.PersonalizationOptions{}, err
	}

	return personalizationOpts, nil
}
