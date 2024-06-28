package internals

import (
	"net/http"
)

func MakeServer() *http.ServeMux {
	mux := http.NewServeMux()
	connManager := NewConnectionManager()
	LoadConfig("config.json")

	mux.HandleFunc("POST /ingest", IngestHandler(connManager))
	return mux
}
