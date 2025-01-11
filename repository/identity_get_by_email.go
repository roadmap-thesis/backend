package repository

import (
	"context"
	"time"

	"github.com/HotPotatoC/roadmap_gen/domain/entity"
	"github.com/HotPotatoC/roadmap_gen/domain/object"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

func (r *Repository) IdentityGetByEmail(ctx context.Context, filter string) (*entity.Identity, error) {
	query, args := psql.Select(
		sm.Columns("id", "name", "email", "password", "created_at", "updated_at"),
		sm.From(entity.IdentityTable),
		sm.Where(psql.Quote("email").EQ(psql.Arg(filter))),
	).MustBuild(ctx)

	var id int
	var name, email, password string
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	identity := &entity.Identity{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  object.NewPasswordFrom(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return identity, nil
}
