package json

import (
	"encoding/json"
	"net/http"

	"github.com/wascript3r/gostr"
)

type jsonError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func newJsonError(err Error) *jsonError {
	return &jsonError{
		Name:    err.Name(),
		Message: gostr.UpperFirst(err.Error()),
	}
}

type jsonRes struct {
	Error *jsonError  `json:"error"`
	Data  interface{} `json:"data"`
}

func SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func NotFound(w http.ResponseWriter, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusNotFound)

	encodeJson(w, NotFoundError, data)
}

func NotFoundCustom(w http.ResponseWriter, err Error, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusNotFound)

	encodeJson(w, err, data)
}

func Forbidden(w http.ResponseWriter, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusForbidden)

	encodeJson(w, ForbiddenError, data)
}

func ForbiddenCustom(w http.ResponseWriter, err Error, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusForbidden)

	encodeJson(w, err, data)
}

func BadRequest(w http.ResponseWriter, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusBadRequest)

	encodeJson(w, BadRequestError, data)
}

func BadRequestCustom(w http.ResponseWriter, err Error, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusBadRequest)

	encodeJson(w, err, data)
}

func InternalError(w http.ResponseWriter, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusInternalServerError)

	encodeJson(w, InternalServerError, data)
}

func InternalErrorCustom(w http.ResponseWriter, err Error, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusInternalServerError)

	encodeJson(w, err, data)
}

func ServeErr(w http.ResponseWriter, err Error, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusOK)

	encodeJson(w, err, data)
}

func ServeJSON(w http.ResponseWriter, data interface{}) {
	SetContentType(w)
	w.WriteHeader(http.StatusOK)

	encodeJson(w, nil, data)
}

func encodeJson(w http.ResponseWriter, err Error, data interface{}) {
	var jsonErr *jsonError
	if err != nil {
		jsonErr = newJsonError(err)
	}

	json.NewEncoder(w).Encode(
		jsonRes{
			Error: jsonErr,
			Data:  data,
		},
	)
}
