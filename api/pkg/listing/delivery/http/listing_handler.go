package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/scraper/api/pkg/listing"
)

type HTTPHandler struct {
	listingUcase listing.Usecase
}

func NewHTTPHandler(r *httprouter.Router, lu listing.Usecase) {
	handler := &HTTPHandler{lu}

	r.POST("/api/listing/register", handler.RegisterListing)
	r.POST("/api/listing/exists", handler.ListingExists)
}

func serveError(w http.ResponseWriter, err error) {
	if err == listing.InvalidInputError {
		httpjson.BadRequestCustom(w, listing.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, listing.UnknownError)
	if code == listing.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) RegisterListing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &listing.RegisterReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.listingUcase.Register(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) ListingExists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &listing.ExistsReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.listingUcase.Exists(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
