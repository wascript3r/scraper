package repository

import "errors"

var ErrTxMismatch = errors.New("transaction type mismatch")

type Transaction interface {
	Commit() error
	Rollback() error
}
