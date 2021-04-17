package query

const (
	DateFormat = "2006-01-02"
)

type QueryRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GetAllRes struct {
	Queries []*QueryRes `json:"queries"`
}

// GetStats

type StatsReq struct {
	ID int `json:"id" validate:"required"`
}

type StatsRes struct {
	URL                      string              `json:"url"`
	Name                     string              `json:"name"`
	Currency                 string              `json:"currency"`
	CurrentAvgPrice          float64             `json:"currentAvgPrice"`
	CurrentRemainingQuantity int                 `json:"currentRemainingQuantity"`
	CurrentAvgSoldPrice      float64             `json:"currentAvgSoldPrice"`
	CurrentSoldQuantity      int                 `json:"currentSoldQuantity"`
	History                  []*HistoryStats     `json:"history"`
	SoldHistory              []*SoldHistoryStats `json:"soldHistory"`
}
