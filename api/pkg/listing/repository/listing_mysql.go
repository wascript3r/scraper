package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	insertMetaSQL     = "INSERT INTO ListingInformation (pkListingId, sellerId, currencyId, title, searchId, conditionTypeId) VALUES (?, ?, ?, ?, ?, ?)"
	insertLocationSQL = "INSERT INTO ListingLocations (fkListingId, locationType, locationId) VALUES (?, ?, ?)"
	insertHistorySQL  = "INSERT INTO ListingHistory (fkListingId, listingPrice, remainingQuantity, dateOfParsing) VALUES (?, ?, ?, ?)"
	existsSQL         = "SELECT EXISTS(SELECT 1 FROM ListingInformation WHERE pkListingId = ?)"
)

type MySQLRepo struct {
	conn *sql.DB
}

func NewMySQLRepo(c *sql.DB) *MySQLRepo {
	return &MySQLRepo{c}
}

func (m *MySQLRepo) NewTx(ctx context.Context) (repository.Transaction, error) {
	return m.conn.BeginTx(ctx, nil)
}

func (m *MySQLRepo) insertMeta(ctx context.Context, q mysql.Querier, ls *domain.ListingMeta) error {
	_, err := q.ExecContext(ctx, insertMetaSQL, ls.ID, ls.SellerID, ls.Currency, ls.Title, ls.SearchQueryID, ls.ConditionID)
	return err
}

func (m *MySQLRepo) InsertMeta(ctx context.Context, ls *domain.ListingMeta) error {
	return m.insertMeta(ctx, m.conn, ls)
}

func (m *MySQLRepo) InsertMetaTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingMeta) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insertMeta(ctx, sqlTx, ls)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) insertLocation(ctx context.Context, q mysql.Querier, ls *domain.ListingLocation) error {
	res, err := q.ExecContext(ctx, insertLocationSQL, ls.ListingID, ls.Type, ls.LocationID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ls.ID = int(id)
	return nil
}

func (m *MySQLRepo) InsertLocation(ctx context.Context, ls *domain.ListingLocation) error {
	return m.insertLocation(ctx, m.conn, ls)
}

func (m *MySQLRepo) InsertLocationTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingLocation) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insertLocation(ctx, sqlTx, ls)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) insertHistory(ctx context.Context, q mysql.Querier, ls *domain.ListingHistory) error {
	res, err := q.ExecContext(ctx, insertHistorySQL, ls.ListingID, ls.Price, ls.RemainingQuantity, ls.ParsedDate)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ls.ID = int(id)
	return nil
}

func (m *MySQLRepo) InsertHistory(ctx context.Context, ls *domain.ListingHistory) error {
	return m.insertHistory(ctx, m.conn, ls)
}

func (m *MySQLRepo) InsertHistoryTx(ctx context.Context, tx repository.Transaction, ls *domain.ListingHistory) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insertHistory(ctx, sqlTx, ls)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) exists(ctx context.Context, q mysql.Querier, id string) (bool, error) {
	var exists bool

	err := q.QueryRowContext(ctx, existsSQL, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *MySQLRepo) Exists(ctx context.Context, id string) (bool, error) {
	return m.exists(ctx, m.conn, id)
}

func (m *MySQLRepo) ExistsTx(ctx context.Context, tx repository.Transaction, id string) (bool, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return false, repository.ErrTxMismatch
	}

	exists, err := m.exists(ctx, sqlTx, id)
	if err != nil {
		sqlTx.Rollback()
		return false, err
	}

	return exists, nil
}
