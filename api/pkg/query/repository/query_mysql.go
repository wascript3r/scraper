package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	insertSQL    = "INSERT INTO SearchRequest (searchUrl, searchExpirityDate, searchName) VALUES(?, ?, ?)"
	getSQL       = "SELECT * FROM SearchRequest WHERE searchId = ?"
	getActiveSQL = "SELECT * FROM SearchRequest WHERE searchExpirityDate > NOW()"
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

func (m *MySQLRepo) insert(ctx context.Context, q mysql.Querier, qs *domain.Query) error {
	res, err := q.ExecContext(ctx, insertSQL, qs.URL, qs.Expiry, qs.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	qs.ID = int(id)
	return nil
}

func (m *MySQLRepo) Insert(ctx context.Context, qs *domain.Query) error {
	return m.insert(ctx, m.conn, qs)
}

func (m *MySQLRepo) InsertTx(ctx context.Context, tx repository.Transaction, qs *domain.Query) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insert(ctx, sqlTx, qs)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) get(ctx context.Context, q mysql.Querier, id int) (*domain.Query, error) {
	qs := &domain.Query{}

	err := q.QueryRowContext(ctx, getSQL, id).Scan(&qs.ID, &qs.URL, &qs.Expiry, &qs.Name)
	if err != nil {
		return nil, mysql.ParseSQLError(err)
	}

	return qs, nil
}

func (m *MySQLRepo) Get(ctx context.Context, id int) (*domain.Query, error) {
	return m.get(ctx, m.conn, id)
}

func (m *MySQLRepo) GetTx(ctx context.Context, tx repository.Transaction, id int) (*domain.Query, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	qs, err := m.get(ctx, sqlTx, id)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return nil, err
	}

	return qs, nil
}

func (m *MySQLRepo) GetActive(ctx context.Context) ([]*domain.Query, error) {
	rows, err := m.conn.QueryContext(ctx, getActiveSQL)
	if err != nil {
		return nil, err
	}

	var queries []*domain.Query
	for rows.Next() {
		qs := &domain.Query{}

		err = rows.Scan(&qs.ID, &qs.URL, &qs.Expiry, &qs.Name)
		if err != nil {
			rows.Close()
			return nil, err
		}

		queries = append(queries, qs)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return queries, nil
}
