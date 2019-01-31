package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func HandleSuccess(w *http.ResponseWriter, result interface{}, logger *logrus.Logger) {
	writer := *w

	marshalled, err := json.Marshal(result)

	if err != nil {
		HandleError(w, 500, "Internal Server Error", "Error marshalling response JSON", err, logger)
		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(marshalled)
}

func HandleError(w *http.ResponseWriter, code int, responseText string, logMessage string, err error, logger *logrus.Logger) {
	errorMessage := ""
	writer := *w

	if err != nil {
		errorMessage = err.Error()
	}

	logger.Error(logMessage, errorMessage)
	writer.WriteHeader(code)
	writer.Write([]byte(responseText))
}
