package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/scraper/api/pkg/query"
)

type Usecase struct {
	queryRepo  query.Repository
	ctxTimeout time.Duration
}

func New(qr query.Repository, t time.Duration) *Usecase {
	return &Usecase{qr, t}
}

func (u *Usecase) GetAll(ctx context.Context) ([]string, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.queryRepo.GetAll(c)
}
