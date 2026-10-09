// Package middleware provides all middlewares for HTTP-server.
package middleware

import (
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// JSON-format for logs
const (
	_infoLevel = 2 // info log level

	_textLogFormat = `time=${time} level=INFO status=${status} method=${method} path=${path}`                       // nolint:lll // output format
	_jsonLogFormat = `{"time":"${time}","level":"INFO","status":"${status}","method":"${method}","path":"${path}"}` // nolint:lll // output format
)

var emptyHandler = func(ctx *fiber.Ctx) error {
	return ctx.Next()
}

type loggerSettings struct {
	level      int
	jsonFormat bool
	locale     string
}

type LoggerOption func(*loggerSettings)

// Logger is a middleware for logging all requests.
func Logger(options ...LoggerOption) fiber.Handler {
	opts := &loggerSettings{
		level:      _infoLevel,
		jsonFormat: false,
		locale:     "UTC",
	}
	for _, option := range options {
		option(opts)
	}

	// not use logger on levels higher than "info"
	if opts.level > _infoLevel {
		return emptyHandler
	}

	// set formatter
	format := _textLogFormat
	if opts.jsonFormat {
		format = _jsonLogFormat
	}

	return logger.New(logger.Config{
		TimeFormat:    time.RFC3339Nano,
		TimeZone:      opts.locale,
		Format:        format + "\n",
		Output:        os.Stderr,
		DisableColors: opts.jsonFormat, // left enabled with text-fmt, disable with JSON-fmt
	})
}

// LoggerWithFormat sets logs format for fiber logger.
// Applies default format (text) if invalid value was given.
// Accepts values: "json", "text".
func LoggerWithFormat(format string) LoggerOption {
	var jsonFormat bool
	if format == "json" {
		jsonFormat = true
	}

	return func(s *loggerSettings) {
		s.jsonFormat = jsonFormat
	}
}

// LoggerWithLevel sets level for fiber logger.
func LoggerWithLevel(level int) LoggerOption {
	return func(s *loggerSettings) {
		s.level = level
	}
}

// LoggerWithLocale sets locale for fiber logger.
func LoggerWithLocale(locale *time.Location) LoggerOption {
	return func(s *loggerSettings) {
		s.locale = locale.String()
	}
}
