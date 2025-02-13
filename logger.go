package log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/sanity-io/litter"

	"github.com/pieterclaerhout/go-formatter"
)

// PrintTimestamp indicates if the log messages should include a timestamp or not
var PrintTimestamp = false

// DebugMode indicates if debug information should be printed or not
var DebugMode = false

// DebugSQLMode indicates if the SQL statements should be logged as debug messages
var DebugSQLMode = false

// TimeZone indicates in which timezone the time should be formatted
var TimeZone *time.Location

// Stdout is the writer to where the stdout messages should be written (defaults to os.Stdout)
var Stdout io.Writer = os.Stdout

// Stderr is the writer to where the stderr messages should be written (defaults to os.Stderr)
var Stderr io.Writer = os.Stderr

// DefaultTimeFormat is the default format to use for the timestamps
var DefaultTimeFormat = "2006-01-02 15:04:05.000"

// TestingTimeFormat is the format to use for the timestamps during testing
var TestingTimeFormat = "test"

// TimeFormat is the format to use for the timestamps
var TimeFormat = DefaultTimeFormat

// OsExit is the function to exit the app when a fatal error happens
var OsExit = os.Exit

// Debug prints a debug message
//
// Only shown if DebugMode is set to true
func Debug(args ...interface{}) {
	if DebugMode {
		message := formatMessage(args...)
		printMessage("DEBUG", message)
	}
}

// DebugSeparator prints a debug separator
//
// Only shown if DebugMode is set to true
func DebugSeparator(args ...interface{}) {
	if DebugMode {
		message := formatMessage(args...)
		message = formatSeparator(message, "=", 80)
		printMessage("DEBUG", message)
	}
}

// DebugSQL formats the SQL statement and prints it as a debug message
//
// Only shown if DebugMode and DebugSQLMode are set to true
func DebugSQL(sql string) {
	if DebugSQLMode {
		message, err := formatter.SQL(sql)
		if err != nil {
			Error(err)
		} else {
			Debug(message)
		}
	}
}

// DebugDump dumps the argument as a debug message with an optional prefix
func DebugDump(arg interface{}, prefix string) {
	message := litter.Sdump(arg)
	if prefix != "" {
		Debug(prefix, message)
	} else {
		Debug(message)
	}
}

// Info prints an info message
func Info(args ...interface{}) {
	message := formatMessage(args...)
	printMessage("INFO ", message)
}

// InfoSeparator prints an info separator
func InfoSeparator(args ...interface{}) {
	message := formatMessage(args...)
	message = formatSeparator(message, "=", 80)
	printMessage("INFO ", message)
}

// InfoDump dumps the argument as an info message with an optional prefix
func InfoDump(arg interface{}, prefix string) {
	message := litter.Sdump(arg)
	if prefix != "" {
		Info(prefix, message)
	} else {
		Info(message)
	}
}

// Warn prints an warning message
func Warn(args ...interface{}) {
	message := formatMessage(args...)
	printMessage("WARN ", message)
}

// WarnDump dumps the argument as a warning message with an optional prefix
func WarnDump(arg interface{}, prefix string) {
	message := litter.Sdump(arg)
	if prefix != "" {
		Warn(prefix, message)
	} else {
		Warn(message)
	}
}

// Error prints an error message to stderr
func Error(args ...interface{}) {
	message := formatMessage(args...)
	printMessage("ERROR", message)
}

// ErrorDump dumps the argument as an err message with an optional prefix to stderr
func ErrorDump(arg interface{}, prefix string) {
	message := litter.Sdump(arg)
	if prefix != "" {
		Error(prefix, message)
	} else {
		Error(message)
	}
}

// StackTrace prints an error message with the stacktrace of err to stderr
func StackTrace(err error) {
	message := formatMessage(FormattedStackTrace(err))
	printMessage("ERROR", message)
}

// FormattedStackTrace returns a formatted stacktrace for err
func FormattedStackTrace(err error) string {
	if cause := causeOfError(err); cause != nil {
		err = cause
	}
	return strings.TrimSpace(errors.Wrap(err, 2).ErrorStack())
}

// Fatal logs a fatal error message to stdout and exits the program with exit code 1
func Fatal(args ...interface{}) {
	message := formatMessage(args...)
	printMessage("FATAL", message)
	OsExit(1)
}

// CheckError checks if the error is not nil and if that's the case, it will print a fatal message and exits the
// program with exit code 1.
//
// If DebugMode is enabled a stack trace will also be printed to stderr
func CheckError(err error) {
	if err != nil {
		printMessage("FATAL", err.Error())
		if DebugMode {
			StackTrace(err)
		}
		OsExit(1)
	}
}
