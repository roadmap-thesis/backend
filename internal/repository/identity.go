package repository

import (
	"context"
	"errors"
	"time"

	"github.com/HotPotatoC/roadmap_gen/internal/domain/entity"
	"github.com/HotPotatoC/roadmap_gen/internal/domain/object"
	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type IdentityRepository struct {
	db DB
}

func NewIdentityRepository(db DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

func (r *IdentityRepository) WithTx(db DB) *IdentityRepository {
	r.db = db
	return r
}

func (r *IdentityRepository) Get(ctx context.Context, col, filter string) (*entity.Identity, error) {
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

func (r *IdentityRepository) GetByID(ctx context.Context, filter string) (*entity.Identity, error) {
	return r.Get(ctx, "id", filter)
}

func (r *IdentityRepository) GetByEmail(ctx context.Context, filter string) (*entity.Identity, error) {
	return r.Get(ctx, "email", filter)
}

func (r *IdentityRepository) Create(ctx context.Context, input *entity.Identity) (*entity.Identity, error) {
	query, args := psql.Insert(
		im.Into(entity.IdentityTable, "name", "email", "password"),
		im.Values(psql.Arg(input.Name, input.Email, input.Password)),
		im.Returning("id", "name", "email", "created_at", "updated_at"),
	).MustBuild(ctx)

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
