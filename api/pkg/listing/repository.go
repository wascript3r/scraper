package listing

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	InsertMeta(ctx context.Context, ls *domain.ListingMeta) error
	InsertMetaTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingMeta) error

	InsertLocation(ctx context.Context, ls *domain.ListingLocation) error
	InsertLocationTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingLocation) error

	InsertHistory(ctx context.Context, ls *domain.ListingHistory) error
	InsertHistoryTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingHistory) error

	InsertSoldHistory(ctx context.Context, ls *domain.ListingSoldHistory) error
	InsertSoldHistoryTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingSoldHistory) error

	Exists(ctx context.Context, id string) (bool, error)
	ExistsTx(ctx context.Context, tx repository.Transaction, id string) (bool, error)
}
