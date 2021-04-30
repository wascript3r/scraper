package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
	"github.com/wascript3r/scraper/api/pkg/query"
)

type HTTPHandler struct {
	queryUcase query.Usecase
}

func NewHTTPHandler(r *httprouter.Router, cors *middleware.Stack, qu query.Usecase) {
	handler := &HTTPHandler{qu}

	r.GET("/api/queries/get", cors.Wrap(handler.GetQueries))
	r.OPTIONS("/api/query/stats", cors.Wrap(handler.GetStats))
	r.POST("/api/query/stats", cors.Wrap(handler.GetStats))
}

func serveError(w http.ResponseWriter, err error) {
	if err == query.InvalidInputError {
		httpjson.BadRequestCustom(w, query.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, query.UnknownError)
	if code == query.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) GetQueries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.queryUcase.GetActive(r.Context())
	if err != nil {
		httpjson.InternalErrorCustom(w, query.UnknownError, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &query.StatsReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.queryUcase.GetStats(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
