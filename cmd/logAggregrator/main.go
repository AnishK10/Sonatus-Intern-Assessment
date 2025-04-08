package main

import (
	"Sonatus-Intern-Assessment/internal/handler"
	"Sonatus-Intern-Assessment/internal/logger"
	"net/http"
)

func main() {
	logger.InitZap()
	defer logger.Zap.Sync()

	http.HandleFunc("/logs", handler.LogHandler)

	logger.Zap.Sugar().Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Zap.Sugar().Fatalw("Server failed", "error", err)
	}
}
