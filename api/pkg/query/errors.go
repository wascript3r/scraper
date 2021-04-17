package query

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	SearchQueryNotFoundError = errcode.New(
		"search_query_not_found",
		errors.New("search query not found"),
	)
)
