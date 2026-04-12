package utils

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type BadRequestError struct {
	Error string `json:"error"`
}

type UnauthorizedError struct {
	Error string `json:"error"`
}

type ForbiddenError struct {
	Error string `json:"error"`
}

type NotFoundError struct {
	Error string `json:"error"`
}

type InternalError struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Error: message,
	})
}

func WriteErrorWithFields(w http.ResponseWriter, status int, message string, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Error:  message,
		Fields: fields,
	})
}

func WriteBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(BadRequestError{Error: message})
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(UnauthorizedError{Error: message})
}

func WriteForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(ForbiddenError{Error: message})
}

func WriteNotFound(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(NotFoundError{Error: message})
}

func WriteInternalError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(InternalError{Error: message})
}

func ValidationError(fields map[string]string) map[string]string {
	return fields
}
