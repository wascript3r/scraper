package query

import "context"

type Usecase interface {
	GetActive(ctx context.Context) (*GetAllRes, error)
}
