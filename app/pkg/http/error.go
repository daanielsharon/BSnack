package httphelper

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

func HandleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		JSONResponse(w, appErr.Code, appErr.Message, nil)
		return
	}

	// Fallback unknown error
	JSONResponse(w, http.StatusInternalServerError, "internal server error", nil)
}
