package repository

import (
	"context"
	"time"

	"github.com/HotPotatoC/roadmap_gen/domain/entity"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

func (r *Repository) IdentityCreate(ctx context.Context, input *entity.Identity) (*entity.Identity, error) {
	query, args := psql.Insert(
		im.Into(entity.IdentityTable, "name", "email", "password"),
		im.Values(psql.Arg(input.Name, input.Email, input.Password)),
		im.Returning("id", "name", "email", "created_at", "updated_at"),
	).MustBuild(ctx)

	log.Info().Msg(query)

	var id int
	var name, email string
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	identity := &entity.Identity{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return identity, nil
}
