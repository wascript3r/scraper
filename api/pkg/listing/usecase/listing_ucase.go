package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/scraper/api/pkg/condition"
	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/listing"
	"github.com/wascript3r/scraper/api/pkg/location"
	"github.com/wascript3r/scraper/api/pkg/photo"
	"github.com/wascript3r/scraper/api/pkg/query"
)

type Usecase struct {
	listingRepo   listing.Repository
	locationRepo  location.Repository
	photoRepo     photo.Repository
	queryRepo     query.Repository
	conditionRepo condition.Repository
	ctxTimeout    time.Duration

	validate listing.Validate
}

func New(lr listing.Repository, lcr location.Repository, pr photo.Repository, qr query.Repository, cr condition.Repository, t time.Duration, v listing.Validate) *Usecase {
	return &Usecase{
		listingRepo:   lr,
		locationRepo:  lcr,
		photoRepo:     pr,
		queryRepo:     qr,
		conditionRepo: cr,
		ctxTimeout:    t,

		validate: v,
	}
}

func (u *Usecase) Register(ctx context.Context, req *listing.RegisterReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return listing.InvalidInputError
	}

	_, err := domain.ToCurrency(req.Currency)
	if err != nil {
		return listing.InvalidCurrencyError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	tx, err := u.listingRepo.NewTx(c)
	if err != nil {
		return err
	}

	exists, err := u.listingRepo.ExistsTx(c, tx, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return listing.AlreadyExistsError
	}

	_, err = u.queryRepo.GetTx(ctx, tx, req.SearchQueryID)
	if err != nil {
		if err == domain.ErrNotFound {
			return listing.SearchQueryNotFoundError
		}
		return err
	}

	_, err = u.conditionRepo.GetTx(ctx, tx, req.Condition)
	if err != nil {
		if err == domain.ErrNotFound {
			return listing.InvalidConditionError
		}
		return err
	}

	// for _, p := range.

	return nil
}
