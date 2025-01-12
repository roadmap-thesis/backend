package provider

import (
	"github.com/HotPotatoC/roadmap_gen/internal/clients"
)

type Provider struct {
	Transaction *TransactionManager
}

func New(clients *clients.Clients) *Provider {
	provider := &Provider{
		Transaction: NewTransactionManager(clients.DB),
	}

	return provider
}
