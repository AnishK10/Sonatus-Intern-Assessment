package handler

import (
	"Sonatus-Intern-Assessment/internal/logger"
	"Sonatus-Intern-Assessment/internal/model"
	"Sonatus-Intern-Assessment/internal/store"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

var logStore = store.NewLogStore()

// Wrapper handler
func LogHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == http.MethodPost {
		IngestLogHandler(w, r)
	} else {
		QueryLogsHandler(w, r)
	}
}

// This will query the logs based on the given params
func QueryLogsHandler(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err1 := time.Parse(time.RFC3339, startStr)
	end, err2 := time.Parse(time.RFC3339, endStr)
	if service == "" || err1 != nil || err2 != nil {
		logger.Zap.Sugar().Warnw("Invalid query parameters", "service", service, "start", startStr, "end", endStr)
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	logs := logStore.GetLogs(service, start, end)
	logger.Zap.Sugar().Infow("Logs retrieved", "service", service, "count", len(logs))
	err := json.NewEncoder(w).Encode(logs)
	if err != nil {
		return
	}
}

// This will ingest the logs in a map where key will be the service name and val will be logs
func IngestLogHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Zap.Sugar().Errorw("Failed to decode request", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry := model.LogEntry{
		Timestamp: req.Timestamp,
		Message:   req.Message,
	}
	logStore.AddLog(req.ServiceName, entry)
	logger.Zap.Sugar().Infow("Log ingested", "service", req.ServiceName, "timestamp", req.Timestamp)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Log added successfully",
	})
}
