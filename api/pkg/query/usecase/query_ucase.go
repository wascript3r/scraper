package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/listing"
	"github.com/wascript3r/scraper/api/pkg/query"
)

type Usecase struct {
	queryRepo  query.Repository
	ctxTimeout time.Duration

	validate query.Validate
}

func New(qr query.Repository, t time.Duration, v query.Validate) *Usecase {
	return &Usecase{qr, t, v}
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

func (u *Usecase) GetStats(ctx context.Context, req *query.StatsReq) (*query.StatsRes, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, listing.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	q, err := u.queryRepo.Get(c, req.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, query.SearchQueryNotFoundError
		}
		return nil, err
	}

	hs, err := u.queryRepo.GetHistoryStats(c, req.ID)
	if err != nil {
		return nil, err
	}

	shs, err := u.queryRepo.GetSoldHistoryStats(c, req.ID)
	if err != nil {
		return nil, err
	}

	qhs := query.ToHistoryStats(hs)
	qshs := query.ToSoldHistoryStats(shs)

	var (
		cPrice        float64
		cQuantity     int
		cSoldPrice    float64
		cSoldQuantity int
	)

	curDate := time.Now().Format(query.DateFormat)

	if len(hs) > 0 && qhs[0].Date == curDate {
		cPrice = qhs[0].AvgPrice
		cQuantity = qhs[0].RemainingQuantity
		qhs = qhs[1:]
	}

	if len(shs) > 0 && qshs[0].Date == curDate {
		cSoldPrice = qshs[0].AvgPrice
		cSoldQuantity = qshs[0].TotalQuantity
		qshs = qshs[1:]
	}

	return &query.StatsRes{
		URL:                      q.URL,
		Name:                     q.Name,
		Currency:                 domain.EURCurrency.String(),
		CurrentAvgPrice:          cPrice,
		CurrentRemainingQuantity: cQuantity,
		CurrentAvgSoldPrice:      cSoldPrice,
		CurrentSoldQuantity:      cSoldQuantity,
		History:                  qhs,
		SoldHistory:              qshs,
	}, nil
}
