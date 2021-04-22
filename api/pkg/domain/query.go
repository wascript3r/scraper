package domain

import "time"

type Query struct {
	ID     int
	URL    string
	Expiry time.Time
	Name   string
}

type QueryHistoryStats struct {
	Date              time.Time
	AvgPrice          float64
	RemainingQuantity int
}

type QuerySoldHistoryStats struct {
	Date          time.Time
	AvgPrice      float64
	TotalQuantity int
}
