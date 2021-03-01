package repository

import (
	"context"
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
	"github.com/wascript3r/scraper/api/pkg/repository"
	"github.com/wascript3r/scraper/api/pkg/repository/mysql"
)

const (
	insertSQL = "INSERT INTO Seller (sellerId, sellerFeedback, satisfactionRate) VALUES (?, ?, ?)"
	getSQL    = "SELECT * FROM Seller WHERE sellerId = ?"
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

func (m *MySQLRepo) insert(ctx context.Context, q mysql.Querier, ss *domain.Seller) error {
	_, err := q.ExecContext(ctx, insertSQL, ss.ID, ss.Feedback, ss.SatisfactionRate)
	return err
}

func (m *MySQLRepo) Insert(ctx context.Context, ss *domain.Seller) error {
	return m.insert(ctx, m.conn, ss)
}

func (m *MySQLRepo) InsertTx(ctx context.Context, tx repository.Transaction, ss *domain.Seller) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return repository.ErrTxMismatch
	}

	err := m.insert(ctx, sqlTx, ss)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	return nil
}

func (m *MySQLRepo) get(ctx context.Context, q mysql.Querier, id string) (*domain.Seller, error) {
	ss := &domain.Seller{}

	err := q.QueryRowContext(ctx, getSQL, id).Scan(&ss.ID, &ss.Feedback, &ss.SatisfactionRate)
	if err != nil {
		return nil, mysql.ParseSQLError(err)
	}

	return ss, nil
}

func (m *MySQLRepo) Get(ctx context.Context, id string) (*domain.Seller, error) {
	return m.get(ctx, m.conn, id)
}

func (m *MySQLRepo) GetTx(ctx context.Context, tx repository.Transaction, id string) (*domain.Seller, error) {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, repository.ErrTxMismatch
	}

	ss, err := m.get(ctx, sqlTx, id)
	if err != nil {
		if err != domain.ErrNotFound {
			sqlTx.Rollback()
		}
		return nil, err
	}

	return ss, nil
}
