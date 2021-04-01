package json

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

type Error interface {
	error
	Name() string
	Original() error
}

var (
	NotFoundError = errcode.New(
		"not_found",
		errors.New("not found"),
	)

	ForbiddenError = errcode.New(
		"forbidden",
		errors.New("forbidden"),
	)

	BadRequestError = errcode.New(
		"bad_request",
		errors.New("bad request"),
	)

	InternalServerError = errcode.New(
		"internal_server_error",
		errors.New("internal server error"),
	)
)
