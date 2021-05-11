package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// Logger wrapper
type Logger struct {
	Log zerolog.Logger
}

var l Logger

func init() {
	var logLevel zerolog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		logLevel = zerolog.DebugLevel
	case "INFO":
		logLevel = zerolog.InfoLevel
	case "WARN":
		logLevel = zerolog.WarnLevel
	case "ERROR":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	l.Log = zerolog.New(os.Stdout)
}

// Infof sends info level event with formatted msg
func Infof(format string, args ...interface{}) {
	l.Log.Info().Msgf(format, args...)
}

// Info sends info level event
func Info(args ...interface{}) {
	l.Log.Info().Msg(fmt.Sprint(args...))
}

// Debugf sends debug level event with formatted msg
func Debugf(format string, args ...interface{}) {
	l.Log.Debug().Msgf(format, args...)
}

// Debug sends debug level event
func Debug(args ...interface{}) {
	l.Log.Debug().Msg(fmt.Sprint(args...))
}

// Errorf sends error level event with formatted msg
func Errorf(format string, args ...interface{}) {
	l.Log.Error().Msgf(format, args...)
}

// Error sends error level event
func Error(args ...interface{}) {
	l.Log.Error().Msg(fmt.Sprint(args...))
}

// Warnf sends warn level event with formatted msg
func Warnf(format string, args ...interface{}) {
	l.Log.Warn().Msgf(format, args...)
}

// Warn sends warn level event
func Warn(args ...interface{}) {
	l.Log.Warn().Msg(fmt.Sprint(args...))
}

// Fatalf sends fatal level event with formatted msg
func Fatalf(format string, args ...interface{}) {
	l.Log.Fatal().Msgf(format, args...)
}

// Fatal sends warn level event
func Fatal(args ...interface{}) {
	l.Log.Fatal().Msg(fmt.Sprint(args...))
}
