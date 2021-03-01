package listing

import "context"

type Usecase interface {
	Register(ctx context.Context, req *RegisterReq) error
	Exists(ctx context.Context, req *ExistsReq) (*ExistsRes, error)
}
