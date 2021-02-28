package location

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ls *domain.Location) error
	InsertTx(ctx context.Context, tx repository.Transaction, ls *domain.Location) error

	Find(ctx context.Context, country, region string) (*domain.Location, error)
	FindTx(ctx context.Context, tx repository.Transaction, country, region string) (*domain.Location, error)
}
