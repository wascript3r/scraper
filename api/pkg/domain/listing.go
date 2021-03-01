package domain

type LocationType int8

const (
	ItemLocationType LocationType = iota
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
