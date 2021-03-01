package listing

type Location struct {
	Country string  `json:"country" validate:"required,lte=100"`
	Region  *string `json:"region" validate:"omitempty,gt=0,lte=100"`
}

type RegisterReq struct {
	ID            string      `json:"id" validate:"required,lte=100"`
	SearchQueryID int         `json:"searchQueryID" validate:"required"`
	Title         string      `json:"title" validate:"required,lte=100"`
	Currency      string      `json:"currency" validate:"required,lte=50"`
	Condition     string      `json:"condition" validate:"required,lte=50"`
	SellerID      string      `json:"sellerID" validate:"required,lte=100"`
	Photos        []string    `json:"photos" validate:"required,dive,url,lte=255"`
	Location      []*Location `json:"location" validate:"required,dive"`
	Shipping      []*Location `json:"shipping" validate:"required,dive"`
}

type ExistsReq struct {
	ID string `json:"id" validate:"required,lte=100"`
}

type ExistsRes struct {
	Exists bool `json:"exists"`
}
