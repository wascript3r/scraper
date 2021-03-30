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

func (u *Usecase) GetActive(ctx context.Context) (*query.GetAllRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	qs, err := u.queryRepo.GetActive(c)
	if err != nil {
		return nil, err
	}

	queries := make([]*query.QueryRes, len(qs))
	for i, q := range qs {
		queries[i] = &query.QueryRes{
			ID:   q.ID,
			Name: q.Name,
			URL:  q.URL,
		}
	}

	return &query.GetAllRes{
		Queries: queries,
	}, nil
}
