package domain

import (
	"fmt"

	"github.com/roadmap-thesis/backend/pkg/apperrors"
)

var (
	ErrNotFound = fmt.Errorf("%w: domain not found", apperrors.NotFound())
)
