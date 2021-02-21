package repository

import (
	"context"
	"database/sql"
)

type MySQLRepo struct {
	conn *sql.DB
}

func NewMySQLRepo() *MySQLRepo {
	return &MySQLRepo{nil}
}

func (m *MySQLRepo) GetAll(ctx context.Context) ([]string, error) {
	return []string{
		"https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=bitcoin&_sacat=0",
		"https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=litecoin&_sacat=0",
		"https://www.ebay.com/sch/i.html?_from=R40&_trksid=p992.m570.l1313&_nkw=dogecoin&_sacat=0",
	}, nil
}
