package photo

import (
	"context"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
)

type Repository interface {
	NewTx(ctx context.Context) (repository.Transaction, error)

	Insert(ctx context.Context, ps *domain.Photo) error
	InsertTx(ctx context.Context, tx repository.Transaction, ps *domain.Photo) error
}
