package models

import "time"

type IngestionRequestBody struct {
	ServiceName string            `json:"service"`
	Timestamp   time.Time         `json:"timestamp"`
	Message     string            `json:"message"`
	Metadata    map[string]string `json:"metadata"`
}
