package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	insertSQL = "INSERT INTO Photos (photoUrl, fkListingId) VALUES(?, ?)"
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

func (m *MySQLRepo) insert(ctx context.Context, q mysql.Querier, ps *domain.Photo) error {
	res, err := q.ExecContext(ctx, insertSQL, ps.URL, ps.ListingID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ps.ID = int(id)
	return nil
}

func (m *MySQLRepo) Insert(ctx context.Context, ps *domain.Photo) error {
	return m.insert(ctx, m.conn, ps)
}

func (m *MySQLRepo) InsertTx(ctx context.Context, tx repository.Transaction, ps *domain.Photo) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insert(ctx, sqlTx, ps)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}
