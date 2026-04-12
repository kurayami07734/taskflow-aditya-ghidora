package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger}
}

func (s *SlogLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &SlogLogEntry{logger: s.logger, request: r}
}

type SlogLogEntry struct {
	logger  *slog.Logger
	request *http.Request
}

func (s *SlogLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	s.logger.Info("request completed",
		"method", s.request.Method,
		"path", s.request.URL.Path,
		"status", status,
		"bytes", bytes,
		"elapsed", elapsed.Seconds(),
	)
}

func (s *SlogLogEntry) Panic(v interface{}, stack []byte) {
	s.logger.Error("panic",
		"panic", v,
		"stack", string(stack),
	)
}

func RecoverWithLog(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovery := recover(); recovery != nil {
					var errMsg string
					switch v := recovery.(type) {
					case error:
						errMsg = v.Error()
					case string:
						errMsg = v
					default:
						errMsg = "unknown panic"
					}

					logger.Error("panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"error", errMsg,
					)

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
