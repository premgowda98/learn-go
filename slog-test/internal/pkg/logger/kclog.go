package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *slog.Logger
	once     sync.Once
)

// Logger levels
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Default configuration constants
const (
	DefaultLogPath    = "logs/app.log"
	DefaultMaxSize    = 5 // megabytes
	DefaultMaxBackups = 3
	DefaultMaxAge     = 28 // days
	DefaultCompress   = true
	DefaultLevel      = LevelInfo
)

// Config holds the configuration for the logger
type Config struct {
	LogPath    string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	Compress   bool
	Level      slog.Level
}

// NewLogger creates a new instance of KCLog with default configuration
func NewLogger(packageName string) *slog.Logger {
	var logger *slog.Logger

	once.Do(func() {
		cfg := Config{
			LogPath:    DefaultLogPath,
			MaxSize:    DefaultMaxSize,
			MaxBackups: DefaultMaxBackups,
			MaxAge:     DefaultMaxAge,
			Compress:   DefaultCompress,
			Level:      DefaultLevel,
		}

		// Ensure log directory exists
		os.MkdirAll(filepath.Dir(cfg.LogPath), 0755)

		// Configure file rotation using lumberjack
		ljLogger := &lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		// Create common handler options
		opts := &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
				}
				return a
			},
		}

		writer := io.MultiWriter(ljLogger, os.Stdout)

		// Create JSON handler with package attribution
		logHandler := slog.NewJSONHandler(writer, opts)
		logger = slog.New(logHandler)

		instance = logger
		slog.SetDefault(instance)
	})

	// Always update the package name, even for existing instance
	return instance.With("package", packageName)
}

// GetLogger returns the singleton logger instance, initializing it if necessary
func GetLogger(packageName string) *slog.Logger {
	if instance == nil {
		return NewLogger(packageName)
	}
	return instance.With("package", packageName)
}
