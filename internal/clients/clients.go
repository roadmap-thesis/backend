package clients

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	OpenAI *OpenAI
	DB     *pgxpool.Pool
}

func New(ctx context.Context) (*Clients, error) {
	c := &Clients{
		OpenAI: NewOpenAIClient(),
	}

	var group errgroup.Group

	group.Go(func() error {
		var err error
		c.DB, err = NewPostgreSQLClient(ctx, config.DatabaseURL())
		if err != nil {
			return errors.Wrap(err, "initializing postgresql")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}
