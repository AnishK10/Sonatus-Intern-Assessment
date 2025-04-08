package test

import (
	"Sonatus-Intern-Assessment/internal/handler"
	"Sonatus-Intern-Assessment/internal/logger"
	"Sonatus-Intern-Assessment/internal/model"
	"Sonatus-Intern-Assessment/internal/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logger.InitZap()
	code := m.Run()
	logger.Zap.Sync()
	os.Exit(code)
}

func TestAddLogAndQueryLog(t *testing.T) {
	logStore := store.NewLogStore()
	timestamp := time.Now().UTC()
	entry := model.LogEntry{
		Timestamp: timestamp,
		Message:   "Unit test log",
	}

	logStore.AddLog("test-service", entry)

	logs := logStore.GetLogs("test-service", timestamp.Add(-1*time.Minute), timestamp.Add(1*time.Minute))
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}
	if logs[0].Message != "Unit test log" {
		t.Errorf("Unexpected log message: %s", logs[0].Message)
	}
}

func TestIngestHandler(t *testing.T) {
	logReq := model.LogRequest{
		ServiceName: "test-service",
		Timestamp:   time.Now().UTC(),
		Message:     "Testing ingest",
	}
	jsonData, _ := json.Marshal(logReq)

	req := httptest.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler.LogHandler(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}
}

func TestQueryHandlerEmptyResult(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?service=unknown&start=2025-01-01T00:00:00Z&end=2025-01-01T01:00:00Z", nil)
	resp := httptest.NewRecorder()

	handler.LogHandler(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	var logs []model.LogEntry
	_ = json.NewDecoder(resp.Body).Decode(&logs)
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(logs))
	}
}
