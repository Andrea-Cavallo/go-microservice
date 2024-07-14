package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// Response is the structure for the HTTP responses
type Response struct {
	Output        interface{}            `json:"output,omitempty"`
	ErrorMessages map[string]interface{} `json:"errorMessages"`
}

// RespondWithJSON writes JSON response to the http.ResponseWriter
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	log := WithContext()
	log.Info("Responding with JSON")

	response := Response{
		Output:        payload,
		ErrorMessages: make(map[string]interface{}), // Initialize as an empty map
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"errorMessages": {"message": "Internal Server Error"}}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}

// RespondWithError writes an error message as JSON response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	log := WithContext()
	log.Info("Responding with error JSON")
	log.Errorf("HTTP %d - %s", code, message)
	response := Response{
		ErrorMessages: map[string]interface{}{"message": message},
	}
	RespondWithJSON(w, code, response)
}

// GenerateUUID generates a new UUID
func GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// EnvOrDefault returns the environment variable if it exists, otherwise the default value
func EnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

// CloseRequestBody closes the request body
func CloseRequestBody(Body io.ReadCloser) {
	log := WithContext()
	log.Info("Closing request body.....")
	err := Body.Close()
	if err != nil {
		log.Errorf("Error closing request body: %v", err)
	}
}
