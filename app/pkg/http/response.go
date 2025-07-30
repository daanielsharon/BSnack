package httphelper

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func getResponseStatus(statusCode int) string {
	if statusCode >= 200 && statusCode < 300 {
		return "success"
	}

	return "failed"
}

func JSONResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Status:  getResponseStatus(statusCode),
		Message: message,
		Data:    data,
	})
}
