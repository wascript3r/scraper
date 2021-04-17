package query

import "github.com/wascript3r/scraper/api/pkg/domain"

type HistoryStats struct {
	Date              string  `json:"date"`
	AvgPrice          float64 `json:"avgPrice"`
	RemainingQuantity int     `json:"remainingQuantity"`
}

func ToHistoryStats(dhs []*domain.QueryHistoryStats) []*HistoryStats {
	hs := make([]*HistoryStats, len(dhs))
	for i, h := range dhs {
		hs[i] = &HistoryStats{
			Date:              h.Date.Format(DateFormat),
			AvgPrice:          h.AvgPrice,
			RemainingQuantity: h.RemainingQuantity,
		}
	}

	return hs
}

type SoldHistoryStats struct {
	Date          string  `json:"date"`
	AvgPrice      float64 `json:"avgPrice"`
	TotalQuantity int     `json:"totalQuantity"`
}

func ToSoldHistoryStats(dhs []*domain.QuerySoldHistoryStats) []*SoldHistoryStats {
	hs := make([]*SoldHistoryStats, len(dhs))
	for i, h := range dhs {
		hs[i] = &SoldHistoryStats{
			Date:          h.Date.Format(DateFormat),
			AvgPrice:      h.AvgPrice,
			TotalQuantity: h.TotalQuantity,
		}
	}

	return hs
}
