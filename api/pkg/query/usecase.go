package query

import "context"

type Usecase interface {
	GetActive(ctx context.Context) (*GetAllRes, error)
	GetStats(ctx context.Context, req *StatsReq) (*StatsRes, error)
}
