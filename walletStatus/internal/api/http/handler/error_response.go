package api

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func makeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := &errorMessage{Message: err.Error()}
	if errE := json.NewEncoder(w).Encode(errResponse); errE != nil {

		return
	}
}
