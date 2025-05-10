package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// Log levels
const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
	LevelWarn  = "WARN"
	LevelDebug = "DEBUG"
)

// Logger represents a custom logger
type Logger struct {
	output io.Writer
}

// New creates a new logger with default output to stdout
func New() *Logger {
	return &Logger{
		output: os.Stdout,
	}
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// getSource returns the file name and line number of the caller
func getSource(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	// Get just the file name, not the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

// formatLog formats the log message according to the required format
func (l *Logger) formatLog(level, msg string, ctx ...interface{}) string {
	// Format: datetime source message<context>
	datetime := time.Now().Format("2006-01-02 15:04:05")
	source := getSource(3) // Skip 3 stack frames to get to the caller
	
	contextStr := ""
	if len(ctx) > 0 {
		parts := make([]string, 0, len(ctx))
		for _, c := range ctx {
			parts = append(parts, fmt.Sprintf("%+v", c))
		}
		contextStr = fmt.Sprintf("<%s>", strings.Join(parts, ", "))
	}
	
	return fmt.Sprintf("%s %s %s %s %s\n", 
		datetime,
		level, 
		source,
		msg,
		contextStr,
	)
}

// Info logs an info message
func (l *Logger) Info(msg string, ctx ...interface{}) {
	fmt.Fprint(l.output, l.formatLog(LevelInfo, msg, ctx...))
}

// Error logs an error message
func (l *Logger) Error(msg string, ctx ...interface{}) {
	fmt.Fprint(l.output, l.formatLog(LevelError, msg, ctx...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, ctx ...interface{}) {
	fmt.Fprint(l.output, l.formatLog(LevelWarn, msg, ctx...))
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, ctx ...interface{}) {
	fmt.Fprint(l.output, l.formatLog(LevelDebug, msg, ctx...))
}