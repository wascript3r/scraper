package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/scraper/api/pkg/auth"
)

var (
	ErrTokenNotProvided = errors.New("auth token is not provided")
	ErrInvalidToken     = errors.New("invalid auth token")
)

type HTTPMiddleware struct {
	authUcase auth.Usecase
}

func NewHTTPMiddleware(au auth.Usecase) *HTTPMiddleware {
	return &HTTPMiddleware{au}
}

func (h *HTTPMiddleware) Authenticated(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token, err := h.ExtractToken(r)
		if err != nil {
			httpjson.Forbidden(w, nil)
			return
		}

		if !h.authUcase.ValidateToken(token) {
			httpjson.Forbidden(w, nil)
			return
		}

		next(w, r, p)
	}
}

func (h *HTTPMiddleware) ExtractToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", ErrTokenNotProvided
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", ErrInvalidToken
	}

	return splitToken[1], nil
}
