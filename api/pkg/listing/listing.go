package listing

type RegisterReq struct {
	ID            string      `json:"id" validate:"required,lte=100"`
	SearchQueryID int         `json:"searchQueryID" validate:"required"`
	Title         string      `json:"title" validate:"required,lte=100"`
	Currency      string      `json:"currency" validate:"required,lte=50"`
	Condition     string      `json:"condition" validate:"required,lte=50"`
	SellerID      string      `json:"sellerID" validate:"required,lte=100"`
	Photos        []string    `json:"photos" validate:"dive,lte=255"`
	Location      []*Location `json:"location" validate:"required"`
	Shipping      []*Location `json:"shipping" validate:"required"`
}

type Location struct {
	Country string `json:"country" validate:"required,lte=100"`
	Region  string `json:"region" validate:"required,lte=100"`
}
