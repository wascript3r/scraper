package domain

import "errors"

type Currency int8

const (
	USDCurrency Currency = iota + 1
	EURCurrency
	GBPCurrency
	AUCurrency
)

func (c Currency) String() string {
	switch c {
	case USDCurrency:
		return "USD"

	case EURCurrency:
		return "EUR"

	case GBPCurrency:
		return "GBP"

	case AUCurrency:
		return "AU"
	}

	return ""
}

var ErrInvalidCurrency = errors.New("invalid currency")

func IsValidCurrency(c Currency) bool {
	switch c {
	case USDCurrency, EURCurrency, GBPCurrency, AUCurrency:
		return true
	}
	return false
}

func ToCurrency(c string) (Currency, error) {
	switch c {
	case "USD":
		return USDCurrency, nil

	case "EUR":
		return EURCurrency, nil

	case "GBP":
		return GBPCurrency, nil

	case "AU":
		return AUCurrency, nil
	}

	return 0, ErrInvalidCurrency
}
