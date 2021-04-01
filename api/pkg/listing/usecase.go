package listing

import "context"

type Usecase interface {
	Register(ctx context.Context, req *RegisterReq) error
	AddHistory(ctx context.Context, req *AddHistoryReq) error
	AddSoldHistory(ctx context.Context, req *AddSoldHistoryReq) error
	Exists(ctx context.Context, req *ExistsReq) (*ExistsRes, error)
}
