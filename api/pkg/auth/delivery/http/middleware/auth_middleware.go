package middleware

import (
	"errors"
	"net/http"
	"strings"

	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/scraper/api/pkg/auth"
)

var (
	ErrTokenNotProvided = errors.New("auth token is not provided")
	ErrInvalidToken     = errors.New("invalid auth token")
)

type HTTPMiddleware struct {
	hnd       http.Handler
	authUcase auth.Usecase
}

func NewHTTPMiddleware(hnd http.Handler, au auth.Usecase) *HTTPMiddleware {
	return &HTTPMiddleware{hnd, au}
}

func (h *HTTPMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.ExtractToken(r)
	if err != nil {
		httpjson.Forbidden(w, nil)
		return
	}

	if !h.authUcase.ValidateToken(token) {
		httpjson.Forbidden(w, nil)
		return
	}

	h.hnd.ServeHTTP(w, r)
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
