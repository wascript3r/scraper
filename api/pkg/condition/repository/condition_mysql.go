package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	getSQL = "SELECT * FROM ItemConditions WHERE conditionName = ?"
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

func (m *MySQLRepo) get(ctx context.Context, q mysql.Querier, name string) (*domain.Condition, error) {
	cs := &domain.Condition{}

	err := q.QueryRowContext(ctx, getSQL, name).Scan(&cs.ID, &cs.Name, &cs.Description)
	if err != nil {
		return nil, mysql.ParseSQLError(err)
	}

	return cs, nil
}

func (m *MySQLRepo) Get(ctx context.Context, name string) (*domain.Condition, error) {
	return m.get(ctx, m.conn, name)
}

func (m *MySQLRepo) GetTx(ctx context.Context, tx repository.Transaction, name string) (*domain.Condition, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	cs, err := m.get(ctx, sqlTx, name)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return nil, err
	}

	return cs, nil
}
