package mcsdk

import (
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(logger *logrus.Logger) *Logger {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Tracef(msg, args...)
}

// Emit a message and key/value pairs at the DEBUG level
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

// Emit a message and key/value pairs at the INFO level
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Emit a message and key/value pairs at the WARN level
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

// Emit a message and key/value pairs at the ERROR level
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (l *Logger) IsTrace() bool {
	return l.logger.Level > logrus.TraceLevel
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (l *Logger) IsDebug() bool {
	return l.logger.Level > logrus.DebugLevel
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (l *Logger) IsInfo() bool {
	return l.logger.Level > logrus.InfoLevel
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (l *Logger) IsWarn() bool {
	return l.logger.Level > logrus.WarnLevel
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (l *Logger) IsError() bool {
	return l.logger.Level > logrus.ErrorLevel
}

// ImpliedArgs returns With key/value pairs
func (l *Logger) ImpliedArgs() []interface{} {
	return nil
}

// Creates a sublogger that will always have the given key/value pairs
func (l *Logger) With(args ...interface{}) hclog.Logger {
	return l
}

// Returns the Name of the logger
func (l *Logger) Name() string {
	return ""
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l *Logger) Named(name string) hclog.Logger {
	return l
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l *Logger) ResetNamed(name string) hclog.Logger {
	return l
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l *Logger) SetLevel(level hclog.Level) {
	logrus.SetLevel(logrus.Level(level))
}

// Return a value that conforms to the stdlib log.Logger interface
func (l *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.Default()
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (l *Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}
