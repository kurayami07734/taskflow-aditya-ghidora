package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log/slog"
)

var defaultLogger *slog.Logger

type rotatingFileHandler struct {
	mu         sync.Mutex
	file       *os.File
	filename   string
	maxSize    int64
	maxBackups int
	h          slog.Handler
}

func (h *rotatingFileHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.h.Enabled(context.Background(), level)
}

func (h *rotatingFileHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.file == nil {
		return nil
	}

	if h.fileSize() >= h.maxSize {
		if err := h.rotate(); err != nil {
			return err
		}
	}

	return h.h.Handle(context.Background(), r)
}

func (h *rotatingFileHandler) fileSize() int64 {
	info, err := h.file.Stat()
	if err != nil {
		return 0
	}
	return info.Size()
}

func (h *rotatingFileHandler) rotate() error {
	h.file.Close()

	for i := h.maxBackups - 1; i >= 1; i-- {
		src := backupName(h.filename, i)
		dst := backupName(h.filename, i+1)
		_ = os.Remove(dst)
		_ = os.Rename(src, dst)
	}

	dst := backupName(h.filename, 1)
	_ = os.Remove(dst)
	_ = os.Rename(h.filename, dst)

	f, err := os.OpenFile(h.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	h.file = f
	h.h = slog.NewJSONHandler(f, nil)
	return nil
}

func backupName(name string, num int) string {
	return name + "." + string(rune('0'+num))
}

func (h *rotatingFileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &rotatingFileHandler{
		mu:         h.mu,
		file:       h.file,
		filename:   h.filename,
		maxSize:    h.maxSize,
		maxBackups: h.maxBackups,
		h:          h.h.WithAttrs(attrs),
	}
}

func (h *rotatingFileHandler) WithGroup(name string) slog.Handler {
	return &rotatingFileHandler{
		mu:         h.mu,
		file:       h.file,
		filename:   h.filename,
		maxSize:    h.maxSize,
		maxBackups: h.maxBackups,
		h:          h.h.WithGroup(name),
	}
}

func newRotatingFileHandler(filename string, maxSize int64, maxBackups int) (*rotatingFileHandler, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &rotatingFileHandler{
		file:       f,
		filename:   filename,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		h:          slog.NewJSONHandler(f, nil),
	}, nil
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(_ context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(context.Background(), level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(_ context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(context.Background(), r); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

func InitLogger(logFile string) *slog.Logger {
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	fh, err := newRotatingFileHandler(logFile, 10<<20, 3)
	if err != nil {
		panic(err)
	}

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	logger := slog.New(&multiHandler{handlers: []slog.Handler{fh, stdoutHandler}})

	slog.SetDefault(logger)
	defaultLogger = logger
	return logger
}

func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

var _ io.Writer = (*nopWriter)(nil)

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func init() {
	_ = time.Now()
}
