package plugins

import (
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

type LogAdapter struct {
	logger zerolog.Logger
	name   string
}

func NewHCLogAdapter(logger zerolog.Logger) *LogAdapter {
	return &LogAdapter{
		logger: logger,
	}
}

func (l *LogAdapter) Log(level hclog.Level, msg string, args ...interface{}) {
	l.logger.WithLevel(l.mapLevel(level)).Msgf(msg, args...)
}

func (l *LogAdapter) Trace(msg string, args ...interface{}) {
	l.logger.Trace().Msgf(msg, args...)
}

// Emit a message and key/value pairs at the DEBUG level
func (l *LogAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

// Emit a message and key/value pairs at the INFO level
func (l *LogAdapter) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

// Emit a message and key/value pairs at the WARN level
func (l *LogAdapter) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

// Emit a message and key/value pairs at the ERROR level
func (l *LogAdapter) Error(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (l *LogAdapter) IsTrace() bool {
	return l.logger.GetLevel() > zerolog.TraceLevel
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (l *LogAdapter) IsDebug() bool {
	return l.logger.GetLevel() > zerolog.DebugLevel
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (l *LogAdapter) IsInfo() bool {
	return l.logger.GetLevel() > zerolog.InfoLevel
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (l *LogAdapter) IsWarn() bool {
	return l.logger.GetLevel() > zerolog.WarnLevel
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (l *LogAdapter) IsError() bool {
	return l.logger.GetLevel() > zerolog.ErrorLevel
}

// ImpliedArgs returns With key/value pairs
func (l *LogAdapter) ImpliedArgs() []interface{} {
	return nil
}

// Creates a sublogger that will always have the given key/value pairs
func (l *LogAdapter) With(args ...interface{}) hclog.Logger {
	return &LogAdapter{
		name:   l.name,
		logger: l.logger.With().Fields(args).Logger(),
	}
}

// Returns the Name of the logger
func (l *LogAdapter) Name() string {
	return l.name
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l *LogAdapter) Named(name string) hclog.Logger {
	fullName := fmt.Sprintf("%s.%s", l.name, name)

	return &LogAdapter{
		name:   fullName,
		logger: l.logger.With().Str("name", fullName).Logger(),
	}
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l *LogAdapter) ResetNamed(name string) hclog.Logger {
	return &LogAdapter{
		name:   name,
		logger: l.logger.With().Str("name", name).Logger(),
	}
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l *LogAdapter) SetLevel(level hclog.Level) {
}

// Return a value that conforms to the stdlib log.LogAdapter interface
func (l *LogAdapter) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.logger, "", 0)
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (l *LogAdapter) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}

func (l *LogAdapter) mapLevel(level hclog.Level) zerolog.Level {
	switch level {
	case hclog.Trace:
		return zerolog.TraceLevel
	case hclog.Debug:
		return zerolog.DebugLevel
	case hclog.Info:
		return zerolog.InfoLevel
	case hclog.Warn:
		return zerolog.WarnLevel
	case hclog.Error:
		return zerolog.ErrorLevel
	case hclog.NoLevel:
		return zerolog.NoLevel
	case hclog.Off:
		return zerolog.Disabled
	default:
		return zerolog.ErrorLevel
	}
}
