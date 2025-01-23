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

type AccountRepository struct {
	db database.Connection
}

func NewAccountRepository(db database.Connection) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) GetByID(ctx context.Context, filter int) (domain.Account, error) {
	accounts, err := r.fetch(ctx, "id", filter)
	if err != nil {
		return domain.Account{}, err
	}

	if len(accounts) == 0 {
		return domain.Account{}, domain.ErrNotFound
	}

	return accounts[0], nil
}

func (r *AccountRepository) GetByEmail(ctx context.Context, filter string) (domain.Account, error) {
	accounts, err := r.fetch(ctx, "email", filter)
	if err != nil {
		return domain.Account{}, err
	}

	if len(accounts) == 0 {
		return domain.Account{}, domain.ErrNotFound
	}

	return accounts[0], nil
}

func (r *AccountRepository) fetch(ctx context.Context, col string, args ...any) ([]domain.Account, error) {
	query, args := psql.Select(
		sm.Columns("id", "name", "email", "password", "created_at", "updated_at"),
		sm.From(domain.AccountTable),
		sm.Where(psql.Quote(col).EQ(psql.Arg(args...))),
	).MustBuild(ctx)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.Password, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) Save(ctx context.Context, input *domain.Account) (domain.Account, error) {
	query, args := psql.Insert(
		im.Into(domain.AccountTable, "name", "email", "password", "created_at", "updated_at"),
		im.Values(psql.Arg(input.Name, input.Email, input.Password, input.CreatedAt, input.UpdatedAt)),
		im.Returning("id", "name", "email", "created_at", "updated_at"),
	).MustBuild(ctx)

	var account domain.Account
	err := r.db.InTx(ctx, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, args...).Scan(
			&account.ID,
			&account.Name,
			&account.Email,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
	})
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}
