package clients

import (
	"context"

	"github.com/pkg/errors"
	"github.com/roadmap-thesis/backend/pkg/config"
	"github.com/roadmap-thesis/backend/pkg/database"
	"github.com/roadmap-thesis/backend/pkg/llm"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	LLM llm.Client
	DB  database.Connection
}

func New(ctx context.Context) (*Clients, error) {
	c := &Clients{
		LLM: llm.NewOpenAiClient(),
		// LLM: llm.NewDeepSeekClient(),
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
