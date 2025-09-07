package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrTooManyReq          = errors.New("too many requests")
	ErrConflict            = errors.New("conflict")
	ErrInternalServerError = errors.New("internal server error")
	ErrServiceUnavailable  = errors.New("service unavailable")
)

type ApiErr interface {
	Status() int
	Error() string
	Detail() string
}

type ApiError struct {
	ErrorStatus  int         `json:"status"`
	ErrorMessage string      `json:"message"`
	ErrorDetail  interface{} `json:"detail"`
}

func (a ApiError) Status() int {
	return a.ErrorStatus
}

func (a ApiError) Error() string {
	return a.ErrorMessage
}

func (a ApiError) Detail() string {
	return fmt.Sprintf("%v", a.ErrorDetail)
}

func NewApiError(status int, message string, errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  status,
		ErrorMessage: message,
		ErrorDetail:  errorDetail,
	}
}

func NewBadRequestError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusBadRequest,
		ErrorMessage: ErrBadRequest.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewNotFoundError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusNotFound,
		ErrorMessage: ErrNotFound.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewUnauthorizedError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusUnauthorized,
		ErrorMessage: ErrUnauthorized.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewForbiddenError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusForbidden,
		ErrorMessage: ErrForbidden.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewTooManyRequestsError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusTooManyRequests,
		ErrorMessage: ErrTooManyReq.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewConflictError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusConflict,
		ErrorMessage: ErrConflict.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewInternalServerError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusInternalServerError,
		ErrorMessage: ErrInternalServerError.Error(),
		ErrorDetail:  errorDetail,
	}
}

func NewServiceUnavailableError(errorDetail interface{}) ApiErr {
	return ApiError{
		ErrorStatus:  http.StatusServiceUnavailable,
		ErrorMessage: ErrServiceUnavailable.Error(),
		ErrorDetail:  errorDetail,
	}
}
