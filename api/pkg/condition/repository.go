package condition

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Get(ctx context.Context, name string) (*domain.Condition, error)
	GetTx(ctx context.Context, tx repository.Transaction, name string) (*domain.Condition, error)
}
