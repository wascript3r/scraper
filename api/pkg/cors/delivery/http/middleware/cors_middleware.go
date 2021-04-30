package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/scraper/api/pkg/auth"
)

type HTTPMiddleware struct {
	authUcase auth.Usecase
}

func NewHTTPMiddleware() *HTTPMiddleware {
	return &HTTPMiddleware{}
}

func (h *HTTPMiddleware) EnableCors(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r, p)
	}
}
