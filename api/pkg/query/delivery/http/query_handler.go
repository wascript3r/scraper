package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/scraper/api/pkg/query"
)

type HTTPHandler struct {
	queryUcase query.Usecase
}

func NewHTTPHandler(r *httprouter.Router, qu query.Usecase) {
	handler := &HTTPHandler{qu}

	r.GET("/api/queries/get", handler.GetQueries)
}

func (h *HTTPHandler) GetQueries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.queryUcase.GetAll(r.Context())
	if err != nil {
		httpjson.InternalErrorCustom(w, query.UnknownError, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}
