package model

import (
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type LogRequest struct {
	ServiceName string    `json:"service_name"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
}
