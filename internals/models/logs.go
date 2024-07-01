package models

import "time"

type LogRecord struct {
	Id        int32     `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Message   string    `json:"message" db:"message"`
}

type LogRecordMetadata struct {
	Id    int32  `json:"id" db:"id"`
	LogId int32  `json:"log_id" db:"log_id"`
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}
