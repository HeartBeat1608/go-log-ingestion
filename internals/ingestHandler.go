package internals

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetServiceNameFromRequest(r *http.Request, body map[string]any) (string, error) {
	svc, ok := body["serviceName"].(string)
	if !ok {
		v := r.Header.Get("X-Service-Name")
		if v == "" {
			return "", fmt.Errorf("no service identifier found")
		}
		svc = v
	}
	return svc, nil
}

func IngestHandler(cm *DBConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			WriteError(w, 500, err)
		}

		svc, ok := body["serviceName"]
		if !ok {
			v := r.Header.Get("X-Service-Name")
			if v == "" {
				WriteError(w, 400, fmt.Errorf("no service identifier found"))
				return
			}
			svc = v
		}

		WriteJSON(w, 200, map[string]interface{}{
			"serviceName": svc,
			"payload":     body,
		})
	}
}
