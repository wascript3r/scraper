package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	insertSQL      = "INSERT INTO Location (country, region) VALUES (?, ?)"
	findNotNullSQL = "SELECT * FROM Location WHERE country = ? AND region = ?"
	findNullSQL    = "SELECT * FROM Location WHERE country = ? AND region IS NULL"
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

func (m *MySQLRepo) insert(ctx context.Context, q mysql.Querier, ls *domain.Location) error {
	res, err := q.ExecContext(ctx, insertSQL, ls.Country, ls.Region)
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

func (m *MySQLRepo) Insert(ctx context.Context, ls *domain.Location) error {
	return m.insert(ctx, m.conn, ls)
}

func (m *MySQLRepo) InsertTx(ctx context.Context, tx repository.Transaction, ls *domain.Location) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insert(ctx, sqlTx, ls)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) find(ctx context.Context, q mysql.Querier, country string, region *string) (*domain.Location, error) {
	var err error
	ls := &domain.Location{}

	if region == nil {
		err = q.QueryRowContext(ctx, findNullSQL, country).Scan(&ls.ID, &ls.Country, &ls.Region)
	} else {
		err = q.QueryRowContext(ctx, findNotNullSQL, country, region).Scan(&ls.ID, &ls.Country, &ls.Region)
	}

	if err != nil {
		return nil, mysql.ParseSQLError(err)
	}

	return ls, nil
}

func (m *MySQLRepo) Find(ctx context.Context, country string, region *string) (*domain.Location, error) {
	return m.find(ctx, m.conn, country, region)
}

func (m *MySQLRepo) FindTx(ctx context.Context, tx repository.Transaction, country string, region *string) (*domain.Location, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ls, err := m.find(ctx, sqlTx, country, region)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return nil, err
	}

	return ls, nil
}
