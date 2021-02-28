package domain

import "time"

type Query struct {
	ID     int
	URL    string
	Expiry time.Time
	Name   string
}
