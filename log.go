package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Zerolog

var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.TraceLevel)

func With() zerolog.Context {
	return Logger.With()
}

func Fatal() *zerolog.Event {
	return Logger.Info()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

// TraceCheck checks if the log level is trace before evaluating the anon fn
func Tracefn(fn func()) {
	if Logger.GetLevel() == zerolog.TraceLevel {
		fn()
	}
}

// Zerolog simple wrappers

// Error for notable errors.
func Errorf(fmts string, a ...interface{}) {
	Logger.Error().Msgf(fmts, a...)
}

// Info for regular messages.
func Infof(fmts string, a ...interface{}) {
	Logger.Info().Msgf(fmts, a...)
}

// Debug for debugging messages.
func Debugf(fmts string, a ...interface{}) {
	Logger.Debug().Msgf(fmts, a...)
}

// Trace for debugging messages.
func Tracef(fmts string, a ...interface{}) {
	Logger.Trace().Msgf(fmts, a...)
}

// Traceln - prints to trace
func Traceln(args ...interface{}) {
	Logger.Trace().Msg(fmt.Sprintln(args...))
}

// Helper functions

// Print - simple print, without timestamp, without regard to log level.
func Print(fmts string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, fmts+"\n", args...)
}

// Log init and settings

// SetLevel sets the log level.
func SetLevel(lvalue zerolog.Level) {
	zerolog.SetGlobalLevel(lvalue)
	jsonLog := os.Getenv("TL_LOG_JSON") == "true"
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(lvalue)
	if !jsonLog {
		// use console logging
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-5s]", i))
		}
		Logger = Logger.Output(output)
	}
	zerolog.DefaultContextLogger = &Logger
	Tracef("Set global log value to %s", lvalue)
}

// setLevelByName sets the log level by string name.
func getLevelByName(lstr string) zerolog.Level {
	switch strings.ToUpper(lstr) {
	case "FATAL":
		return zerolog.FatalLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "INFO":
		return zerolog.InfoLevel
	case "DEBUG":
		return zerolog.DebugLevel
	case "TRACE":
		return zerolog.TraceLevel
	}
	return zerolog.InfoLevel
}

func init() {
	SetLevel(getLevelByName(os.Getenv("TL_LOG")))
}
