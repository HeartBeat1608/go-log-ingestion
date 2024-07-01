package internals

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/HeartBeat1608/logmaster/internals/models"
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
		log.Printf("> %s", r.URL.Path)

		var body models.IngestionRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			WriteError(w, 500, err)
			return
		}

		// if service name is not in the body, we will check the header
		if body.ServiceName == "" {
			v := r.Header.Get("X-Service-Name")
			if v == "" {
				WriteError(w, 400, fmt.Errorf("no service identifier found"))
				return
			}
			body.ServiceName = v
		}

		db, err := cm.GetConnection(body.ServiceName)
		if err != nil {
			WriteError(w, 500, err)
			return
		}

		tx, err := db.BeginTxx(context.Background(), nil)
		if err != nil {
			WriteError(w, 500, err)
			return
		}
		res, err := tx.Exec("INSERT INTO logs (timestamp, message) VALUES (?, ?)", body.Timestamp, body.Message)
		if err != nil {
			WriteError(w, 500, err)
			return
		}
		rec_id, err := res.LastInsertId()
		if err != nil {
			WriteError(w, 500, err)
			return
		}

		for key, value := range body.Metadata {
			if _, err = tx.Exec("INSERT INTO metadata (log_id, key, value) VALUES (?, ?, ?)", rec_id, key, value); err != nil {
				WriteError(w, 500, err)
				return
			}
		}

		if err = tx.Commit(); err != nil {
			WriteError(w, 500, err)
			return
		}

		WriteJSON(w, 200, map[string]any{
			"payload": body,
		})
	}
}
