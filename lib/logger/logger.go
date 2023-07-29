package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	SetLogger(NewLogger(&LoggerConfig{
		Level:  DebugLevel,
		Output: "stdout",
		File:   "log.txt",
	}))
}

// Logger is an interface for log output
type Logger interface {
	// Debug output
	Debug(v ...interface{})
	// Debugf output
	Debugf(format string, v ...interface{})
	// Info output
	Info(v ...interface{})
	// Infof output
	Infof(format string, v ...interface{})
	// Warn output
	Warn(v ...interface{})
	// Warnf output
	Warnf(format string, v ...interface{})
	// Error output
	Error(v ...interface{})
	// Errorf output
	Errorf(format string, v ...interface{})
	// Fatal output
	Fatal(v ...interface{})
	// Fatalf output
	Fatalf(format string, v ...interface{})
}

type logger struct {
	level  Level
	output *log.Logger
}

var loggerInstance Logger

// Level is the log level
type Level int

const (
	// DebugLevel is the debug level
	DebugLevel Level = iota
	// InfoLevel is the info level
	InfoLevel
	// WarnLevel is the warn level
	WarnLevel
	// ErrorLevel is the error level
	ErrorLevel
	// FatalLevel is the fatal level
	FatalLevel
)

// LoggerConfig is the config for logger
type LoggerConfig struct {
	// Log level
	Level Level
	// Log output
	Output string
	// Log file
	File string
}

// NewLogger returns a new logger
func NewLogger(config *LoggerConfig) Logger {
	file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	output := log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags)

	return &logger{
		level:  config.Level,
		output: output,
	}
}

// SetLogger sets the logger
func SetLogger(logger Logger) {
	loggerInstance = logger
}

// Debug output
func Debug(v ...interface{}) {
	loggerInstance.Debug(v...)
}

// Debugf output
func Debugf(format string, v ...interface{}) {
	loggerInstance.Debugf(format, v...)
}

// Info output
func Info(v ...interface{}) {
	loggerInstance.Info(v...)
}

// Infof output
func Infof(format string, v ...interface{}) {
	loggerInstance.Infof(format, v...)
}

// Warn output
func Warn(v ...interface{}) {
	loggerInstance.Warn(v...)
}

// Warnf output
func Warnf(format string, v ...interface{}) {
	loggerInstance.Warnf(format, v...)
}

// Error output
func Error(v ...interface{}) {
	loggerInstance.Error(v...)
}

// Errorf output
func Errorf(format string, v ...interface{}) {
	loggerInstance.Errorf(format, v...)
}

// Fatal output
func Fatal(v ...interface{}) {
	loggerInstance.Fatal(v...)
}

// Fatalf output
func Fatalf(format string, v ...interface{}) {
	loggerInstance.Fatalf(format, v...)
}

// Debug output
func (l *logger) Debug(v ...interface{}) {
	if l.level <= DebugLevel {
		l.output.Output(2, fmt.Sprint(v...))
	}
}

// Debugf output
func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.output.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info output
func (l *logger) Info(v ...interface{}) {
	if l.level <= InfoLevel {
		l.output.Output(2, fmt.Sprint(v...))
	}
}

// Infof output
func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.output.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn output
func (l *logger) Warn(v ...interface{}) {
	if l.level <= WarnLevel {
		l.output.Output(2, fmt.Sprint(v...))
	}
}

// Warnf output
func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.output.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error output
func (l *logger) Error(v ...interface{}) {
	if l.level <= ErrorLevel {
		l.output.Output(2, fmt.Sprint(v...))
	}
}

// Errorf output
func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.output.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal output
func (l *logger) Fatal(v ...interface{}) {
	if l.level <= FatalLevel {
		l.output.Output(2, fmt.Sprint(v...))
		os.Exit(1)
	}
}

// Fatalf output
func (l *logger) Fatalf(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		l.output.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}
