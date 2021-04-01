package seller

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ss *domain.Seller) error
	InsertTx(ctx context.Context, tx repository.Transaction, ss *domain.Seller) error

	Get(ctx context.Context, id string) (*domain.Seller, error)
	GetTx(ctx context.Context, tx repository.Transaction, id string) (*domain.Seller, error)
}
