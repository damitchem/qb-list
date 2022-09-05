package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type Logger struct {
	minLevel LogLevel
}

func (l *Logger) Log(level LogLevel, message string, details map[string]interface{}) {
	if level < l.minLevel {
		return
	}

	log := make(map[string]interface{})
	log["level"] = levelNames[level]
	now := time.Now().UTC()
	log["timestamp"] = now.String()
	log["message"] = message
	log["details"] = details
	output, _ := json.Marshal(log)
	fmt.Println(string(output))
}

func New(opts ...func(*Logger)) *Logger {
	logger := &Logger{
		minLevel: Info,
	}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

func MinLogLevel(level LogLevel) func(*Logger) {
	return func(logger *Logger) {
		logger.minLevel = level
	}
}

func (l *Logger) Trace(message string, details ...interface{}) {
	l.Log(Trace, message, mapDetails(details))
}

func (l *Logger) Debug(message string, details ...interface{}) {
	l.Log(Debug, message, mapDetails(details))
}

func (l *Logger) Info(message string, details ...interface{}) {
	l.Log(Info, message, mapDetails(details))
}

func (l *Logger) Warn(message string, details ...interface{}) {
	l.Log(Warn, message, mapDetails(details))
}

func (l *Logger) Error(message string, details ...interface{}) {
	l.Log(Error, message, mapDetails(details))
}

func (l *Logger) Critical(message string, details ...interface{}) {
	l.Log(Critical, message, mapDetails(details))
}

func mapDetails(details []interface{}) map[string]interface{} {
	l := len(details)
	detailsMap := make(map[string]interface{})
	if l == 0 {
		return detailsMap
	}
	if l%2 != 0 {
		details = append(details, "ERROR VALUE")
	}
	for i, detail := range details {
		if i%2 != 0 {
			continue
		}
		val := details[i+1]
		detailsMap[detail.(string)] = val
	}
	return detailsMap
}
