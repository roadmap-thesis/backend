package repository

import (
	"context"
	"errors"
	"time"

	"github.com/HotPotatoC/roadmap_gen/internal/domain/entity"
	"github.com/HotPotatoC/roadmap_gen/internal/domain/object"
	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

func (r *Repository) IdentityGet(ctx context.Context, col, filter string) (*entity.Identity, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	identity := &entity.Identity{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  object.Password(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return identity, nil
}

func (r *Repository) IdentityGetByID(ctx context.Context, filter string) (*entity.Identity, error) {
	return r.IdentityGet(ctx, "id", filter)
}

func (r *Repository) IdentityGetByEmail(ctx context.Context, filter string) (*entity.Identity, error) {
	return r.IdentityGet(ctx, "email", filter)
}
