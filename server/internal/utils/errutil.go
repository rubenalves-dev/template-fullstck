package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("conflict")
)

func MapError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "RESOURCE_NOT_FOUND"
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "UNAUTHORIZED"
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN"
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, "BAD_REQUEST"
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, "CONFLICT"
	default:
		return http.StatusInternalServerError, "INTERNAL_SERVER_ERROR"
	}
}
