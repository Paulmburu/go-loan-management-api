package response

import (
	"errors"
	"go-loan-management-api/internal/apperror"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrValidation):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperror.ErrNotFound):
		Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, apperror.ErrConflict):
		Error(w, http.StatusConflict, err.Error())
	default:
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
