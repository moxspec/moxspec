package loglet

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type level int

// These are the different logging levels
const (
	DEBUG level = iota
	INFO
	WARN
	ERROR
	FATAL
	NONE
)

const (
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	white        = 37
	lightGray    = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	lightWhite   = 97
)

var (
	minLv  level
	writer io.Writer
	term   bool
)

func init() {
	minLv = DEBUG
	writer = os.Stderr
	term = isTerminal(writer)
}

// Logger represents a logger
type Logger struct {
	header string
}

// NewLogger creates and initializes a Logger
func NewLogger(header string) *Logger {
	l := new(Logger)
	l.header = header
	return l
}

// SetLevel sets the logger level
func SetLevel(lv level) {
	if DEBUG <= lv && lv <= NONE {
		minLv = lv
	}
}

// SetOutput sets the logger output
func SetOutput(w io.Writer) {
	writer = w
	term = isTerminal(w)
}

func isTerminal(w io.Writer) bool {
	switch w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(w.(*os.File).Fd()))
	}
	return false
}

func (l Logger) output(lv level, format string, a ...interface{}) {
	if minLv > lv {
		return
	}

	var lvStr string
	var color int
	switch lv {
	case DEBUG:
		lvStr = "debug"
		color = cyan
	case INFO:
		lvStr = "info"
		color = green
	case WARN:
		lvStr = "warn"
		color = yellow
	case ERROR:
		lvStr = "error"
		color = red
	case FATAL:
		lvStr = "fatal"
		color = red
	default:
		lvStr = "unknown"
	}

	line := fmt.Sprintf("[%s][%s] ", lvStr, l.header) + format
	if term {
		line = fmt.Sprintf("\x1b[%dm", color) + line + "\x1b[0m\n"
	} else {
		line = line + "\n"
	}

	if len(a) > 0 {
		fmt.Fprintf(writer, line, a...)
	} else {
		fmt.Fprintf(writer, line)
	}
}

// Debug logs a message at level Debug
func (l Logger) Debug(a interface{}) {
	l.Debugf(convertString(a))
}

// Debugf logs a message at level Debug
func (l Logger) Debugf(format string, a ...interface{}) {
	l.output(DEBUG, format, a...)
}

// Info logs a message at level Info
func (l Logger) Info(a interface{}) {
	l.Infof(convertString(a))
}

// Infof logs a message at level Info
func (l Logger) Infof(format string, a ...interface{}) {
	l.output(INFO, format, a...)
}

// Warn logs a message at level Warn
func (l Logger) Warn(a interface{}) {
	l.Warnf(convertString(a))
}

// Warnf logs a message at level Warn
func (l Logger) Warnf(format string, a ...interface{}) {
	l.output(WARN, format, a...)
}

// Error logs a message at level Error
func (l Logger) Error(a interface{}) {
	l.Errorf(convertString(a))
}

// Errorf logs a message at level Error
func (l Logger) Errorf(format string, a ...interface{}) {
	l.output(ERROR, format, a...)
}

// Fatal logs a message at level Fatal and panic program
func (l Logger) Fatal(a interface{}) {
	l.Fatalf(convertString(a))
}

// Fatalf logs a message at level Fatal and panic program
func (l Logger) Fatalf(format string, a ...interface{}) {
	l.output(FATAL, format, a...)
	panic("abort")
}

func convertString(a interface{}) string {
	switch a.(type) {
	case string:
		return a.(string)
	case fmt.Stringer:
		return a.(fmt.Stringer).String()
	}
	return fmt.Sprintf("%+v", a)
}
