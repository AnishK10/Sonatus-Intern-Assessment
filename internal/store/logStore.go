package store

import (
	"Sonatus-Intern-Assessment/internal/logger"
	"Sonatus-Intern-Assessment/internal/model"
	"sort"
	"sync"
	"time"
)

const cronTime = 1 * time.Minute

type LogStore struct {
	mu   sync.RWMutex
	logs map[string][]model.LogEntry
}

func NewLogStore() *LogStore {
	store := &LogStore{logs: make(map[string][]model.LogEntry)}
	go store.cleanupExpiredLogs()
	return store
}

func (ls *LogStore) AddLog(serviceName string, entry model.LogEntry) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.logs[serviceName] = append(ls.logs[serviceName], entry)
	sort.Slice(ls.logs[serviceName], func(i, j int) bool {
		return ls.logs[serviceName][i].Timestamp.Before(ls.logs[serviceName][j].Timestamp)
	})
}

func (ls *LogStore) GetLogs(serviceName string, start, end time.Time) []model.LogEntry {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	var logs []model.LogEntry
	for _, log := range ls.logs[serviceName] {
		if !log.Timestamp.Before(start) && !log.Timestamp.After(end) {
			logs = append(logs, log)
		}
	}
	return logs
}

func (ls *LogStore) cleanupExpiredLogs() {
	for {
		time.Sleep(cronTime) // Potential improvement making it configurable
		expiration := time.Now().Add(-1 * time.Hour)
		ls.mu.Lock()
		for svc, entries := range ls.logs {
			filtered := entries[:0]
			for _, entry := range entries {
				if entry.Timestamp.After(expiration) {
					filtered = append(filtered, entry)
				}
			}
			ls.logs[svc] = filtered
		}
		ls.mu.Unlock()
		logger.Zap.Sugar().Debug("Expired logs cleaned up")
	}
}
