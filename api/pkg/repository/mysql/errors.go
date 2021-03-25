package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/wascript3r/scraper/api/pkg/domain"
)

type ErrCode uint16

const (
	DuplicateEntryErrCode ErrCode = 1062
)

func ParseSQLError(err error) error {
	switch err {
	case sql.ErrNoRows:
		return domain.ErrNotFound
	}
	return err
}

func ParseMySQLError(err error) error {
	if e, ok := err.(*mysql.MySQLError); ok {
		switch ErrCode(e.Number) {
		case DuplicateEntryErrCode:
			return domain.ErrExists
		}
	}
	return err
}
