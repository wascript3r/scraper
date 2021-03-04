package listing

import (
	"errors"

	"github.com/wascript3r/cryptopay/pkg/errcode"
)

var (
	// Error codes

	InvalidInputError = errcode.InvalidInputError
	UnknownError      = errcode.UnknownError

	InvalidCurrencyError = errcode.New(
		"invalid_currency_error",
		errors.New("invalid currency error"),
	)

	AlreadyExistsError = errcode.New(
		"listing_already_exists",
		errors.New("listing already exists"),
	)

	SearchQueryNotFoundError = errcode.New(
		"search_query_not_found",
		errors.New("search query not found"),
	)

	InvalidConditionError = errcode.New(
		"invalid_condition",
		errors.New("invalid condition"),
	)
)
