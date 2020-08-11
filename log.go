package log

import (
	"os"
)

type (
	// Level of log
	Level uint
)

func (lv Level) String() string {
	switch lv {
	case LevelFatal:
		return "Fatal"
	case LevelError:
		return "Error"
	case LevelWarning:
		return "Warning"
	case LevelInfo:
		return "Info"
	case LevelDebug:
		return "Debug"
	}

	return "UNKNOWN"
}

func (lv Level) Color() int {
	switch lv {
	case LevelFatal:
		return 35
	case LevelError:
		return 31
	case LevelWarning:
		return 33
	case LevelInfo:
		return 32
	case LevelDebug:
		return 36
	}

	return 36
}

const (
	// LevelFatal .
	LevelFatal Level = iota
	// LevelError .
	LevelError
	// LevelWarning .
	LevelWarning
	// LevelInfo .
	LevelInfo
	// LevelDebug .
	LevelDebug
)

var builtin *Logger // the builtin Logger

func init() {
	builtin, _ = NewLogger()
}

// Fatal .
func Fatal(args ...interface{}) {
	builtin.Fatal(args...)
	os.Exit(1)
}

// Fatalf .
func Fatalf(format string, args ...interface{}) {
	builtin.Fatalf(format, args...)
	os.Exit(1)
}

// Error .
func Error(args ...interface{}) {
	builtin.Error(args...)
}

// Errorf .
func Errorf(format string, args ...interface{}) {
	builtin.Errorf(format, args...)
}

// Warn .
func Warn(args ...interface{}) {
	builtin.Warn(args...)
}

// Warnf .
func Warnf(format string, args ...interface{}) {
	builtin.Warnf(format, args...)
}

// Info .
func Info(args ...interface{}) {
	builtin.Info(args...)
}

// Infof .
func Infof(format string, args ...interface{}) {
	builtin.Infof(format, args...)
}

// Debug .
func Debug(args ...interface{}) {
	builtin.Debug(args...)
}

// Debugf .
func Debugf(format string, args ...interface{}) {
	builtin.Debugf(format, args...)
}

// WithField .
func WithField(key string, value interface{}) *entry {
	return builtin.WithField(key, value)
}

// WithFields .
func WithFields(fields Fields) *entry {
	return builtin.WithFields(fields)
}

// SetLogLevel .
func SetLogLevel(level Level) {
	builtin.SetLogLevel(level)
}

func SetCallerReporter(b bool) {
	builtin.SetCallerReporter(b)
}
