package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, status int, err error) {
	res := errorResponse{Message: err.Error()}
	w.WriteHeader(status)

	if werr := json.NewEncoder(w).Encode(res); werr != nil {
		slog.Error("failed to write response", "error", werr.Error())
	}
}

func SendOkResponse[T any](w http.ResponseWriter, res T) {
	w.WriteHeader(http.StatusOK)
	if werr := json.NewEncoder(w).Encode(res); werr != nil {
		slog.Error("failed to write response", "error", werr.Error())
	}
}
