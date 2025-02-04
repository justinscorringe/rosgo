package ros

import (
	"fmt"
	"log"
	"os"
)

//LogLevel is the global int to set the log level severity of a logger
type LogLevel int

const (
	//LogLevelDebug is zero value of LogLevel, most verbose
	LogLevelDebug LogLevel = iota
	//LogLevelInfo is default value of LogLevel
	LogLevelInfo
	//LogLevelWarn is a warning LogLevel
	LogLevelWarn
	//LogLevelError is a standard Error LogLevel
	LogLevelError
	//LogLevelFatal is a fatal Error LogLevel
	LogLevelFatal
)

//Logger interface functions for log level severities
type Logger interface {
	Severity() LogLevel
	SetSeverity(severity LogLevel)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type defaultLogger struct {
	severity LogLevel
}

//NewDefaultLogger returns a new defaultLogger with severity LogLevelInfo
func NewDefaultLogger() Logger {
	logger := new(defaultLogger)
	logger.severity = LogLevelInfo
	return logger
}

func (logger *defaultLogger) Severity() LogLevel {
	return logger.severity
}

func (logger *defaultLogger) SetSeverity(severity LogLevel) {
	logger.severity = severity
}

func (logger *defaultLogger) Debug(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelDebug) {
		msg := fmt.Sprintf("[DEBUG] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Debugf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelDebug) {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (logger *defaultLogger) Info(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		msg := fmt.Sprintf("[INFO] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Infof(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		log.Printf("[INFO] "+format, v...)
	}
}

func (logger *defaultLogger) Warn(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		msg := fmt.Sprintf("[WARN] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Warnf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		log.Printf("[WARN] "+format, v...)
	}
}

func (logger *defaultLogger) Error(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		msg := fmt.Sprintf("[ERROR] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Errorf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		log.Printf("[ERROR]"+format, v...)
	}
}

func (logger *defaultLogger) Fatal(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		msg := fmt.Sprintf("[FATAL] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

func (logger *defaultLogger) Fatalf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		log.Printf("[FATAL] "+format, v...)
		os.Exit(1)
	}
}
