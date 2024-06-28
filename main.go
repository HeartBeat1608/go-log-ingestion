package main

import (
	"log"
	"net/http"

	"github.com/HeartBeat1608/logmaster/internals"
)

var (
	listenAddr string = ":5000"
)

func main() {
	server := internals.MakeServer()

	log.Fatal(http.ListenAndServe(listenAddr, server))
}
