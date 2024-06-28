package internals

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	data := map[string]interface{}{
		"error": err.Error(),
	}
	WriteJSON(w, statusCode, data)
}
