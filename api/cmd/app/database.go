package main

import (
	"database/sql"
)

func openDatabase(driver, connStr string) (*sql.DB, error) {
	conn, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
