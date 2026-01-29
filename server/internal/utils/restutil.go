package domain

import (
	"encoding/json"
	"net/http"
)

type ResponseEnvelope struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func RenderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ResponseEnvelope{
		Data: data,
	})
}

func RenderError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ResponseEnvelope{
		Error: &ErrorDetail{
			Code:    code,
			Message: msg,
		},
	})
}
