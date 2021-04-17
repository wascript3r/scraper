package query

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, qs *domain.Query) error
	InsertTx(ctx context.Context, tx repository.Transaction, qs *domain.Query) error

	Get(ctx context.Context, id int) (*domain.Query, error)
	GetTx(ctx context.Context, tx repository.Transaction, id int) (*domain.Query, error)

	GetActive(ctx context.Context) ([]*domain.Query, error)
	GetHistoryStats(ctx context.Context, id int) ([]*domain.QueryHistoryStats, error)
	GetSoldHistoryStats(ctx context.Context, id int) ([]*domain.QuerySoldHistoryStats, error)
}
