package main

import (
	"log"
	"net/http"

	"github.com/HeartBeat1608/logmaster/internals"
	"github.com/HeartBeat1608/logmaster/views"
)

var (
	listenAddr string = ":9956"
)

func main() {
	_ = internals.LoadConfig("config.json")
	server := MakeServer()

	log.Printf("Listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, server))
}

func MakeServer() *http.ServeMux {
	mux := http.NewServeMux()
	connManager := internals.NewConnectionManager()
	vm := views.NewViewManager(connManager)

	mux.HandleFunc("GET /", internals.RootHandler(connManager))
	mux.HandleFunc("GET /health", internals.HealthHandler)
	mux.HandleFunc("POST /ingest", internals.IngestHandler(connManager))

	mux.HandleFunc("GET /{service}/view/all", vm.RenderAllLogs)
	return mux
}
