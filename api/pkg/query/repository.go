package query

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]string, error)
}
