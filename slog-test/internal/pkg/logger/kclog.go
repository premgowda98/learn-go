package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *slog.Logger
	once     sync.Once
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

const (
	DefaultLogDir     = "logs"
	DefaultMaxSize    = 5
	DefaultMaxBackups = 3
	DefaultMaxAge     = 28
	DefaultCompress   = true
	DefaultLevel      = LevelDebug
)

type Config struct {
	LogDir     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      slog.Level
}

func NewLogger(packageName string) *slog.Logger {
	var logger *slog.Logger

	//singleton
	once.Do(func() {
		cfg := Config{
			LogDir:     DefaultLogDir,
			MaxSize:    DefaultMaxSize,
			MaxBackups: DefaultMaxBackups,
			MaxAge:     DefaultMaxAge,
			Compress:   DefaultCompress,
			Level:      DefaultLevel,
		}

		os.MkdirAll(cfg.LogDir, 0755)

		errorLogger := &lumberjack.Logger{
			Filename:   filepath.Join(cfg.LogDir, "error.log"),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		infoLogger := &lumberjack.Logger{
			Filename:   filepath.Join(cfg.LogDir, "info.log"),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		debugLogger := &lumberjack.Logger{
			Filename:   filepath.Join(cfg.LogDir, "debug.log"),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

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

		handler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(&leveledHandler{
			defaultHandler: handler,
			errorWriter:    errorLogger,
			infoWriter:     infoLogger,
			debugWriter:    debugLogger,
			opts:           opts,
		})

		instance = logger
		slog.SetDefault(instance)
	})

	return instance
}

type CustomHandler struct {
	slog.Handler
	w    io.Writer
	opts slog.HandlerOptions
}

func NewCustomHandler(w io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		Handler: slog.NewTextHandler(w, opts),
		w:       w,
		opts:    *opts,
	}
}

func getCaller(skip int) *slog.Source {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	wd, err := os.Getwd()
	if err == nil {
		if relPath, err := filepath.Rel(wd, file); err == nil {
			file = relPath
		}
	}

	return &slog.Source{
		File: file,
		Line: line,
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	level := strings.ToUpper(r.Level.String())
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")

	var otherOpts string = ""
	// fmt.Println(r.sou)
	r.Attrs(func(a slog.Attr) bool {
		otherOpts += a.Key
		otherOpts += " "
		otherOpts += a.Value.String()
		return true
	})

	source := getCaller(2)

	output := fmt.Sprintf("%s%s %s %s %s\n",
		level[:1],
		timeStr,
		slog.String("source", fmt.Sprintf("%s:%d", source.File, source.Line)),
		r.Message,
		otherOpts,
	)

	_, err := h.w.Write([]byte(output))
	return err
}

type leveledHandler struct {
	defaultHandler slog.Handler
	errorWriter    io.Writer
	infoWriter     io.Writer
	debugWriter    io.Writer
	opts           *slog.HandlerOptions
}

func (h *leveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.defaultHandler.Enabled(ctx, level)
}

func (h *leveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &leveledHandler{
		defaultHandler: h.defaultHandler.WithAttrs(attrs),
		errorWriter:    h.errorWriter,
		infoWriter:     h.infoWriter,
		debugWriter:    h.debugWriter,
		opts:           h.opts,
	}
}

func (h *leveledHandler) WithGroup(name string) slog.Handler {
	return &leveledHandler{
		defaultHandler: h.defaultHandler.WithGroup(name),
		errorWriter:    h.errorWriter,
		infoWriter:     h.infoWriter,
		debugWriter:    h.debugWriter,
		opts:           h.opts,
	}
}

func (h *leveledHandler) Handle(ctx context.Context, r slog.Record) error {
	var handler slog.Handler
	switch {
	case r.Level >= LevelError:
		handler = NewCustomHandler(h.errorWriter, h.opts)
	case r.Level >= LevelInfo:
		handler = NewCustomHandler(h.infoWriter, h.opts)
	default:
		handler = NewCustomHandler(h.debugWriter, h.opts)
	}

	return handler.Handle(ctx, r)
}
