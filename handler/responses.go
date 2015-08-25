package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jeffbmartinez/log"
)

/*
WriteJSONResponse writes a json response (with correct http header).
*/
func WriteJSONResponse(response http.ResponseWriter, message interface{}, statusCode int) {
	responseString, err := json.Marshal(message)
	if err != nil {
		log.Errorf("Couldn't marshal json: %v", err)

		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(""))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write([]byte(responseString))
}
