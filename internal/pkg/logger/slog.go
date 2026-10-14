// Package logger provides init function to setup global slog logger.
package logger

import (
	"io"
	"log/slog"
	"os"
	"time"
)

const (
	_debugLevel  = 1 // debug log level
	_infoLevel   = 2 // info log level
	_warnLevel   = 3 // warn log level
	_errorLevel  = 4 // error log level
	_silentLevel = 5 // silent log level
)

// slogParams represents params for default slog logger.
type slogParams struct {
	// JSON logs format if true else text (by default)
	jsonFormat bool
	// locale (time zone) for log messages
	locale *time.Location
	// Logs level (1 - debug, 2 - info, 3 - warn, 4 - error, 5 - silent)
	level int
	// Output writer
	out io.Writer
}

// SlogOption represents a param for slog logger.
type SlogOption func(*slogParams)

// InitSlog sets up default slog logger for application with given options.
// Also it sets default log logger (standart lib).
// Default logger uses text format and info log level.
// Options can be set with "With..." funcs.
func InitSlog(options ...SlogOption) {
	// default options
	params := &slogParams{
		jsonFormat: false,
		locale:     time.UTC,
		level:      _infoLevel,
		out:        os.Stderr,
	}
	// apply custom options
	for _, option := range options {
		option(params)
	}

	if params.level == _silentLevel {
		devNullLog := slog.New(slog.NewJSONHandler(io.Discard, nil))
		slog.SetDefault(devNullLog)
		return
	}

	// options for slog with log level
	handlerOpts := &slog.HandlerOptions{
		Level: slog.Level(translateSlogLogLevel(params.level)),
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			// use time in locale
			if attr.Key == "time" {
				if t, ok := attr.Value.Any().(time.Time); ok {
					attr.Value = slog.TimeValue(t.In(params.locale))
				}
			}
			return attr
		},
	}

	// create new slog instance with specified format
	var newLog *slog.Logger
	if params.jsonFormat {
		newLog = slog.New(slog.NewJSONHandler(params.out, handlerOpts))
	} else {
		newLog = slog.New(slog.NewTextHandler(params.out, handlerOpts))
	}
	slog.SetDefault(newLog)
}

// WithFormat sets format for slog logger.
// Applies default format (text) if invalid value was given.
// Accepts values: "json", "text".
func WithFormat(format string) SlogOption {
	var jsonFormat bool
	if format == "json" {
		jsonFormat = true
	}

	return func(s *slogParams) {
		s.jsonFormat = jsonFormat
	}
}

// WithJSONFormat sets JSON format for slog logger.
func WithJSONFormat() SlogOption {
	return func(s *slogParams) {
		s.jsonFormat = true
	}
}

// WithLocale sets locale for slog logger.
func WithLocale(locale *time.Location) SlogOption {
	return func(s *slogParams) {
		s.locale = locale
	}
}

// WithDebugLevel sets debug log level for slog logger.
func WithDebugLevel() SlogOption {
	return func(s *slogParams) {
		s.level = _debugLevel
	}
}

// WithInfoLevel sets info log level for slog logger.
func WithInfoLevel() SlogOption {
	return func(s *slogParams) {
		s.level = _infoLevel
	}
}

// WithWarnLevel sets warn log level for slog logger.
func WithWarnLevel() SlogOption {
	return func(s *slogParams) {
		s.level = _warnLevel
	}
}

// WithErrorLevel sets error log level for slog logger.
func WithErrorLevel() SlogOption {
	return func(s *slogParams) {
		s.level = _errorLevel
	}
}

// WithSilentLevel sets silent log level for slog logger.
func WithSilentLevel() SlogOption {
	return func(s *slogParams) {
		s.level = _silentLevel
	}
}

// WithLogLevel sets log level for slog logger.
// Appplies default level (2 - info) if invalid value was given.
// Accepts values: 1 - debug, 2 - info, 3 - warn, 4 - error, 5 - silent.
func WithLogLevel(level int) SlogOption {
	if level < _debugLevel || level > _silentLevel {
		level = _infoLevel
	}
	return func(s *slogParams) {
		s.level = level
	}
}

// translateSlogLogLevel translates sequently
// log level (1..2..3..4) to slog level (-4..0..4..8).
func translateSlogLogLevel(logLevel int) int {
	return logLevel*4 - 8 // nolint:mnd // sequence formula
}
