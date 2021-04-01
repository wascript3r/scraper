package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("no items found")
	ErrNullValue           = errors.New("item value is null")
	ErrExists              = errors.New("item already exists")
	ErrInvalidItem         = errors.New("invalid item")
	ErrInvalidParamInput   = errors.New("invalid input parameter")
)
