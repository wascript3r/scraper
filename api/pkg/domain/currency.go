package domain

import "errors"

type Currency int8

const (
	USDCurrency Currency = iota
	EURCurrency
)

var ErrInvalidCurrency = errors.New("invalid currency")

func IsValidCurrency(c Currency) bool {
	switch c {
	case USDCurrency, EURCurrency:
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
	}
	return 0, ErrInvalidCurrency
}
