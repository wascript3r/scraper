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
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/seller"
)

type Usecase struct {
	listingRepo   listing.Repository
	locationRepo  location.Repository
	photoRepo     photo.Repository
	queryRepo     query.Repository
	conditionRepo condition.Repository
	sellerRepo    seller.Repository
	ctxTimeout    time.Duration

	hasher   listing.Hasher
	validate listing.Validate
}

func New(lr listing.Repository, lcr location.Repository, pr photo.Repository, qr query.Repository, cr condition.Repository,
	sr seller.Repository, t time.Duration, h listing.Hasher, v listing.Validate) *Usecase {
	return &Usecase{
		listingRepo:   lr,
		locationRepo:  lcr,
		photoRepo:     pr,
		queryRepo:     qr,
		conditionRepo: cr,
		sellerRepo:    sr,
		ctxTimeout:    t,

		hasher:   h,
		validate: v,
	}
}

func (u *Usecase) Register(ctx context.Context, req *listing.RegisterReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return listing.InvalidInputError
	}

	curr, err := domain.ToCurrency(req.Currency)
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

	_, err = u.queryRepo.GetTx(c, tx, req.SearchQueryID)
	if err != nil {
		if err == domain.ErrNotFound {
			return listing.SearchQueryNotFoundError
		}
		return err
	}

	cond, err := u.conditionRepo.GetTx(c, tx, req.Condition)
	if err != nil {
		if err == domain.ErrNotFound {
			return listing.InvalidConditionError
		}
		return err
	}

	_, err = u.sellerRepo.GetTx(c, tx, req.SellerID)
	if err != nil {
		if err != domain.ErrNotFound {
			return err
		}

		ss := &domain.Seller{ID: req.SellerID}
		err = u.sellerRepo.InsertTx(c, tx, ss)
		if err != nil {
			return err
		}
	}

	meta := &domain.ListingMeta{
		ID:            req.ID,
		SellerID:      req.SellerID,
		Currency:      curr,
		Title:         req.Title,
		SearchQueryID: req.SearchQueryID,
		ConditionID:   cond.ID,
	}

	err = u.listingRepo.InsertMetaTx(c, tx, meta)
	if err != nil {
		return err
	}

	err = u.insertPhotos(c, tx, meta.ID, req.Photos)
	if err != nil {
		return err
	}

	err = u.insertLocations(c, tx, meta.ID, req.Location, domain.ItemLocationType)
	if err != nil {
		return err
	}

	err = u.insertLocations(c, tx, meta.ID, req.Shipping, domain.ShippingLocationType)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) insertPhotos(ctx context.Context, tx repository.Transaction, listingID string, photos []string) error {
	for _, p := range photos {
		err := u.photoRepo.InsertTx(ctx, tx, &domain.Photo{
			URL:       p,
			ListingID: listingID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Usecase) insertLocations(ctx context.Context, tx repository.Transaction, listingID string, locations []*listing.Location, t domain.LocationType) error {
	for _, l := range locations {
		ls, err := u.locationRepo.FindTx(ctx, tx, l.Country, l.Region)
		if err != nil {
			if err != domain.ErrNotFound {
				return err
			}

			ls = &domain.Location{
				Country: l.Country,
				Region:  l.Region,
			}

			err = u.locationRepo.InsertTx(ctx, tx, ls)
			if err != nil {
				return err
			}
		}

		err = u.listingRepo.InsertLocationTx(ctx, tx, &domain.ListingLocation{
			ListingID:  listingID,
			Type:       t,
			LocationID: ls.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Usecase) AddHistory(ctx context.Context, req *listing.AddHistoryReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return listing.InvalidInputError
	}

	date, err := time.Parse(u.validate.GetDateTimeFormat(), req.ParsedDate)
	if err != nil {
		return listing.CannotParseDateError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	exists, err := u.listingRepo.Exists(c, req.ListingID)
	if err != nil {
		return err
	}
	if !exists {
		return listing.ListingNotFoundError
	}

	history := &domain.ListingHistory{
		ListingID:         req.ListingID,
		Price:             req.Price,
		RemainingQuantity: req.RemainingQuantity,
		ParsedDate:        date,
	}

	return u.listingRepo.InsertHistory(c, history)
}

func (u *Usecase) AddSoldHistory(ctx context.Context, req *listing.AddSoldHistoryReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return listing.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	exists, err := u.listingRepo.Exists(c, req.ListingID)
	if err != nil {
		return err
	}
	if !exists {
		return listing.ListingNotFoundError
	}

	tx, err := u.listingRepo.NewTx(c)
	if err != nil {
		return err
	}

	for _, h := range req.History {
		date, err := time.Parse(u.validate.GetDateTimeFormat(), h.Date)
		if err != nil {
			return listing.CannotParseDateError
		}

		history := &domain.ListingSoldHistory{
			IDHash:        u.hasher.HashSoldRecord(h),
			ListingID:     req.ListingID,
			Price:         h.Price,
			Quantity:      h.Quantity,
			PurchasedDate: date,
		}

		err = u.listingRepo.InsertSoldHistoryTx(c, tx, history)
		if err != nil && err != domain.ErrExists {
			return err
		}
	}

	return tx.Commit()
}

func (u *Usecase) Exists(ctx context.Context, req *listing.ExistsReq) (*listing.ExistsRes, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, listing.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	exists, err := u.listingRepo.Exists(c, req.ID)
	if err != nil {
		return nil, err
	}

	return &listing.ExistsRes{
		Exists: exists,
	}, nil
}
