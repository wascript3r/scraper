package mysql

import (
	"database/sql"

	"github.com/wascript3r/scraper/api/pkg/domain"
)

type ErrCode string

const (
	UniqueViolationErrCode ErrCode = "23505"
	CheckViolationErrCode  ErrCode = "23514"
)

func ParseSQLError(err error) error {
	switch err {
	case sql.ErrNoRows:
		return domain.ErrNotFound
	}
	return err
}
