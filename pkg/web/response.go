package web

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {

	// Establecer el encabezado Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Establecer el código de estado del response
	w.WriteHeader(statusCode)

	// Serializar la respuesta a JSON y enviarla en la respuesta
	err := json.NewEncoder(w).Encode(Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})

	if err != nil {
		Error(w, http.StatusInternalServerError, "Error al serializar la respuesta a JSON")
		return
	}

}

func Error(w http.ResponseWriter, statusCode int, message string) {

	// Establecer el encabezado Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Establecer el código de estado del response
	w.WriteHeader(statusCode)

	// Serializar la respuesta a JSON y enviarla en la respuesta
	json.NewEncoder(w).Encode(Response{
		Code:    statusCode,
		Message: message,
	})

}
