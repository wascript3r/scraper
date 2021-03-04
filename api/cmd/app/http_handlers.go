package main

import (
	"errors"
	"net/http"

	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
)

var (
	MethodNotAllowedError = errcode.New(
		"method_not_allowed",
		errors.New("method not allowed"),
	)
	MethodNotAllowedHnd = http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			httpjson.BadRequestCustom(w, MethodNotAllowedError, nil)
		},
	)

	NotFoundHnd = http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			httpjson.NotFound(w, nil)
		},
	)
)
