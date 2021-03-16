package domain

import "time"

type LocationType int8

const (
	ItemLocationType LocationType = iota + 1
	ShippingLocationType
)

type ListingMeta struct {
	ID            string
	SellerID      string
	Currency      Currency
	Title         string
	SearchQueryID int
	ConditionID   int
}

type ListingLocation struct {
	ID         int
	ListingID  string
	Type       LocationType
	LocationID int
}

type ListingHistory struct {
	ID                int
	ListingID         string
	Price             float64
	RemainingQuantity int
	ParsedDate        time.Time
}
