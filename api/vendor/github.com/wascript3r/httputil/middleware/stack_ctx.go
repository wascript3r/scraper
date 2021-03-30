package middleware

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/httputil"
)

type MiddlewareCtx func(next httputil.HandleCtx) httputil.HandleCtx

type StackCtx struct {
	middlewares []MiddlewareCtx
}

func NewCtx() *StackCtx {
	return &StackCtx{nil}
}

func (s *StackCtx) Use(m MiddlewareCtx) {
	s.middlewares = append(s.middlewares, m)
}

func (s *StackCtx) Wrap(ctx context.Context, fn httputil.HandleCtx) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l := len(s.middlewares)
		if l == 0 {
			fn(ctx, w, r, p)
			return
		}

		var result httputil.HandleCtx
		result = s.middlewares[l-1](fn)

		for i := l - 2; i >= 0; i-- {
			result = s.middlewares[i](result)
		}

		result(ctx, w, r, p)
	}
}
