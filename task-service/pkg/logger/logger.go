package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Level represents log severity
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Entry represents a structured log entry compatible with GCP Cloud Logging
type Entry struct {
	Severity  Level  `json:"severity"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId,omitempty"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
	Duration  string `json:"duration,omitempty"`
}

var output = log.New(os.Stdout, "", 0)

// write outputs a structured JSON log entry
func write(entry Entry) {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(entry)
	output.Println(string(data))
}

// Info logs an informational message
func Info(msg string) {
	write(Entry{Severity: LevelInfo, Message: msg})
}

// Warn logs a warning message
func Warn(msg string) {
	write(Entry{Severity: LevelWarn, Message: msg})
}

// Error logs an error message
func Error(msg string) {
	write(Entry{Severity: LevelError, Message: msg})
}

// Request logs an HTTP request with structured fields
func Request(requestID, method, path string, status int, duration time.Duration) {
	write(Entry{
		Severity:  LevelInfo,
		Message:   "request completed",
		RequestID: requestID,
		Method:    method,
		Path:      path,
		Status:    status,
		Duration:  duration.String(),
	})
}
