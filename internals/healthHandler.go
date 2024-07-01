package internals

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, 200, map[string]interface{}{
		"status": "running",
	})
}
