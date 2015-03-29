// Package logger provides a custom logging abstract over the standard out logging of golang.
// All logging should by go to stdout according to 12-factor principles.
// Logging levels are based on RFC 5424 - http://www.rfc-base.org/rfc-5424.html#
package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Standard labels
const (

	//  RFC 5424 log levels.
	Emergency = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug
)

const (
	// ANSI 8 colours.
	foregroundBlack = iota + 30
	foregroundRed
	foregroundGreen
	foregroundYellow
	foregroundBlue
	foregroundMagenta
	foregroundCyan
	foregroundLightGrey
	_
	foregroundDefault

	colourFormat = "[\x1b[%dm%s\x1b[0m] "
)

var (
	// Log labels.
	labels = []string{"[EMERGENCY] ",
		"[ALERT] ",
		"[CRITICAL] ",
		"[ERROR] ",
		"[WARNING] ",
		"[NOTICE] ",
		"[INFO] ",
		"[DEBUG] ",
	}
)

// Logger provides a datastructure for all logging state.
type Logger struct {
	logger *log.Logger
	level  int
	labels []string
}

// New is a factory method to return a new logger instance.
func New(level int, colours bool) *Logger {
	flags := log.Lshortfile | log.Ldate | log.Lmicroseconds
	pre := fmt.Sprintf("[%d] ", os.Getpid())
	if level < 0 {
		level = Info
	}
	l := &Logger{
		logger: log.New(os.Stdout, pre, flags),
		level:  level,
	}

	if colours {
		l.SetColouredLabels()
	} else {
		l.SetPlainLabels()
	}
	return l
}

// SetLogLevel allows a user to set the log level of the logger
func (l *Logger) SetLogLevel(logLevel int) error {

	if logLevel < Emergency || logLevel > Debug {
		return errors.New(fmt.Sprintf("%d log level arg is not in valid range.", logLevel))
	}
	l.level = logLevel
	return nil
}

// GetLogLevel returns the current log level of the logger
func (l *Logger) GetLogLevel() int {
	return l.level
}

// SetPlainLabels sets the message labels to simple text output.
func (l *Logger) SetPlainLabels() {
	copy(l.labels, labels)
}

// SetColouredLabels sets the message labels to colourized text output.
func (l *Logger) SetColouredLabels() {
	l.labels = make([]string, 0)
	for i, lbl := range labels {
		var colour int
		switch i {
		case Emergency, Alert, Critical, Error:
			colour = foregroundRed
		case Warning:
			colour = foregroundYellow
		case Notice:
			colour = foregroundGreen
		case Debug:
			colour = foregroundBlue
		default:
			colour = foregroundDefault
		}
		l.labels = append(l.labels, fmt.Sprintf(colourFormat, colour, lbl))
	}
}

// Emergencyf prints an emergency message to the system log,
// This is considered an unrecoverable error and the application also exits, unless dont exit = true.
func (l *Logger) Emergencyf(exit bool, format string, v ...interface{}) {
	if l.level >= Emergency {
		l.Output(3, labels[Emergency], format, v...)
	}
	if exit == true {
		os.Exit(1)
	}
}

// Alertf prints an alert message to the system log.
func (l *Logger) Alertf(format string, v ...interface{}) {
	if l.level >= Alert {
		l.Output(3, labels[Alert], format, v...)
	}
}

// Criticalf prints a critical message to the system log.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	if l.level >= Critical {
		l.Output(3, labels[Critical], format, v...)
	}
}

// Errorf prints an error message to the system log.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level >= Error {
		l.Output(3, labels[Error], format, v...)
	}
}

// Warningf prints a warning message to the system log.
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.level >= Warning {
		l.Output(3, labels[Warning], format, v...)
	}
}

// Noticef prints a notice message to the system log.
func (l *Logger) Noticef(format string, v ...interface{}) {
	if l.level >= Notice {
		l.Output(3, labels[Notice], format, v...)
	}
}

// Infof prints an informational message to the system log.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level >= Info {
		l.Output(3, labels[Info], format, v...)
	}
}

// Debugf prints a debug message to the system log.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level >= Debug {
		l.Output(3, labels[Debug], format, v...)
	}
}

// output prints a message directly into the system log. Normally, you should use level message functions.
// so that level can trap the write.
func (l *Logger) Output(calldepth int, label string, format string, v ...interface{}) error {
	var d int = 2
	if calldepth > 0 {
		d = calldepth
	}
	return l.logger.Output(d, fmt.Sprintf(label+format, v...))
}
