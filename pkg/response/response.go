package response

import (
	"errors"
	"net/http"
	"strings"

	apiError "github.com/ClearingHouse/pkg/api_error"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorResponseBuilder(err error) (int, ErrorResponse) {
	var apiErr apiError.ApiErr

	switch {
	case errors.As(err, &apiErr):
		return apiErr.Status(), ErrorResponse{
			Status:  apiErr.Status(),
			Message: apiErr.Error(),
			Error:   apiErr.Detail(),
		}

	case strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "Field validation"):
		// Gin binding errors often contain "validation" in the error string
		return http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Validation error",
			Error:   err.Error(),
		}

	default:
		return http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: apiError.ErrInternalServerError.Error(),
			Error:   "-",
		}
	}
}
