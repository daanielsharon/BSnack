package validation

import (
	httphelper "bsnack/app/pkg/http"
	"net/http"
)

func ValidateEnoughPoint(customerPoints int, pointRequired int) error {
	if customerPoints < pointRequired {
		return httphelper.NewAppError(http.StatusBadRequest, "Not enough points")
	}
	return nil
}
