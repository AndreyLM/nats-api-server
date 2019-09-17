package apiserver

import (
	"encoding/json"
	"net/http"
)

func makeJSONResponse(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response to json", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
