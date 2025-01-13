package repository

import (
	"context"
	"errors"
	"time"

	"github.com/HotPotatoC/roadmap_gen/internal/domain"
	"github.com/HotPotatoC/roadmap_gen/internal/domain/object"
	"github.com/HotPotatoC/roadmap_gen/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type AccountRepository interface {
	Get(ctx context.Context, col string, filter any) (*domain.Account, error)
	GetByID(ctx context.Context, filter int) (*domain.Account, error)
	GetByEmail(ctx context.Context, filter string) (*domain.Account, error)

	Create(ctx context.Context, input *domain.Account) (*domain.Account, error)
}

type accountRepository struct {
	db database.Connection
}

func NewAccountRepository(db database.Connection) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Get(ctx context.Context, col string, filter any) (*domain.Account, error) {
	query, args := psql.Select(
		sm.Columns("id", "name", "email", "password", "created_at", "updated_at"),
		sm.From(domain.AccountTable),
		sm.Where(psql.Quote(col).EQ(psql.Arg(filter))),
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

	account := &domain.Account{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  object.Password(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return account, nil
}

func (r *accountRepository) GetByID(ctx context.Context, filter int) (*domain.Account, error) {
	return r.Get(ctx, "id", filter)
}

func (r *accountRepository) GetByEmail(ctx context.Context, filter string) (*domain.Account, error) {
	return r.Get(ctx, "email", filter)
}

func (r *accountRepository) Create(ctx context.Context, input *domain.Account) (*domain.Account, error) {
	query, args := psql.Insert(
		im.Into(domain.AccountTable, "name", "email", "password"),
		im.Values(psql.Arg(input.Name, input.Email, input.Password)),
		im.Returning("id", "name", "email", "created_at", "updated_at"),
	).MustBuild(ctx)

	var id int
	var name, email string
	var createdAt, updatedAt time.Time

	err := r.db.InTx(ctx, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &createdAt, &updatedAt)
	})
	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return account, nil
}
