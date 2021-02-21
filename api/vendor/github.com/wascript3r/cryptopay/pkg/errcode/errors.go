package errcode

import "errors"

const (
	UnknownErrName ErrName = "unknown"
)

var (
	InvalidInputError = New(
		"invalid_input",
		errors.New("input could not pass the validations"),
	)

	UnknownError = New(
		UnknownErrName,
		errors.New("unknown error"),
	)
)
