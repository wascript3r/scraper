package query

import "context"

type Usecase interface {
	GetAll(ctx context.Context) ([]string, error)
}
