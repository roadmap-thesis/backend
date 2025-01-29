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

type TopicRepository struct {
	db database.Connection
}

func NewTopicRepository(db database.Connection) *TopicRepository {
	return &TopicRepository{
		db: db,
	}
}

func (r *TopicRepository) GetBySlug(ctx context.Context, slug string) (domain.Topic, error) {
	ctx, span := tracer.Start(ctx, "(*TopicRepository.GetBySlug)", trace.WithAttributes(attribute.String("slug", slug)))
	defer span.End()

	topics, err := r.fetch(ctx, "slug", slug)
	if err != nil {
		return domain.Topic{}, err
	}

	if len(topics) == 0 {
		return domain.Topic{}, domain.ErrTopicNotFound
	}

	return topics[0], nil
}

func (r *TopicRepository) fetch(ctx context.Context, col string, args ...any) ([]domain.Topic, error) {
	ctx, span := tracer.Start(ctx, "(*TopicRepository.fetch)", trace.WithAttributes(attribute.String("col", col)))
	defer span.End()

	query, args := psql.Select(
		sm.Columns(
			psql.Quote(domain.TopicTable, "id"),
			psql.Quote(domain.TopicTable, "roadmap_id"),
			psql.F("COALESCE", psql.Quote(domain.TopicTable, "parent_id"), 0),
			psql.Quote(domain.TopicTable, "title"),
			psql.Quote(domain.TopicTable, "slug"),
			psql.Quote(domain.TopicTable, "description"),
			psql.Quote(domain.TopicTable, "order"),
			psql.Quote(domain.TopicTable, "finished"),
			psql.Quote(domain.TopicTable, "created_at"),
			psql.Quote(domain.TopicTable, "updated_at"),
		),
		sm.From(domain.TopicTable),
		sm.Where(psql.Quote(domain.TopicTable, col).EQ(psql.Arg(args...))),
	).MustBuild(ctx)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var topic domain.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.RoadmapID,
			&topic.ParentID,
			&topic.Title,
			&topic.Slug,
			&topic.Description,
			&topic.Order,
			&topic.Finished,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}
