package clients

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/internal/config"
	"github.com/HotPotatoC/roadmap_gen/internal/database"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	OpenAI *OpenAI
	DB     database.Connection
}

func New(ctx context.Context) (*Clients, error) {
	c := &Clients{
		OpenAI: NewOpenAIClient(),
	}

	var group errgroup.Group

	group.Go(func() error {
		var err error
		c.DB, err = database.New(ctx, config.DatabaseURL())
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

func (c *Clients) Close() {
	c.DB.Close()
	log.Info().Msg("clients shutdown complete")
}
