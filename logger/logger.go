/**
V2 Logger that is simalar to Spring Boot or NestJS logger supporting colours and names for log levels.
*/

package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type (
	LogMessage struct {
		AppName string
		Time    string
		Name    string
		Color   string
		Level   string
		Msg     string
	}

	ILogStream interface {
		Write(msg LogMessage)
	}

	LoggerOption struct {
		// Name of the logger. This will be used as the logger's name to identify the logger
		Name string

		// If true, the logger will not print anything to stdout
		Silent bool

		// Extra streams to write to
		ExtraStreams []ILogStream
	}

	Logger interface {
		Trace(msg string)
		Debug(msg string)
		Log(msg string)
		Warn(msg string)
		Error(msg string)
		Fatal(msg string)
		Panic(msg string)

		// f functions
		Tracef(format string, args ...interface{})
		Debugf(format string, args ...interface{})
		Logf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		Panicf(format string, args ...interface{})
	}

	LoggerImpl struct {
		Logger
		Streams []ILogStream
		name    string
	}
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"

	// Global prefix for the logger. Default is "GoApp"
	globalPrefix = "GoApp"

	// Global extra streams
	globalExtraStreams []ILogStream
)

func NewLogger(o LoggerOption) Logger {
	var ss []ILogStream

	if !o.Silent {
		ss = append(ss, NewStdOutStream())
	}

	// Append global extra streams
	ss = append(ss, globalExtraStreams...)

	// Append extra streams
	ss = append(ss, o.ExtraStreams...)

	return &LoggerImpl{
		name:    o.Name,
		Streams: ss,
	}
}

// Set the global prefix for the logger. If set, this name will be added to all log messages as a prefix.
func SetAppName(prefix string) {
	globalPrefix = prefix
}

// Add the streams to the global extra streams. This will be added to all loggers.
func AddGlobalExtraStream(streams []ILogStream) {
	globalExtraStreams = append(globalExtraStreams, streams...)
}

func (l *LoggerImpl) logWithColor(color, logLevel, msg string) {
	// Get the current time here so that all streams have the same time
	ct := time.Now().Format(time.RFC3339)

	// For all streams
	for _, stream := range l.Streams {
		stream.Write(LogMessage{
			AppName: globalPrefix,
			Name:    l.name,
			Time:    ct,
			Color:   color,
			Level:   logLevel,
			Msg:     msg,
		})
	}

}

func (l *LoggerImpl) logWithColorf(color, logLevel, format string, args ...interface{}) {
	l.logWithColor(color, logLevel, fmt.Sprintf(format, args...))
}

func (l *LoggerImpl) Trace(msg string) {
	l.logWithColor(Cyan, "TRACE", msg)
}

func (l *LoggerImpl) Debug(msg string) {
	// Check if runtime is development
	if os.Getenv("RUNTIME") != "development" {
		return
	}

	// Blue
	l.logWithColor(Blue, "DEBUG", msg)
}

func (l *LoggerImpl) Log(msg string) {
	// Green
	l.logWithColor(Green, "LOG", msg)
}

func (l *LoggerImpl) Warn(msg string) {
	// Yellow
	l.logWithColor(Yellow, "WARN", msg)
}

func (l *LoggerImpl) Error(msg string) {
	// Red
	l.logWithColor(Red, "ERROR", msg)
}

func (l *LoggerImpl) Fatal(msg string) {
	// Red + Fatal + Exit(1)
	l.logWithColor(Red, "FATAL", msg)
	log.Fatal(msg)
}

func (l *LoggerImpl) Panic(msg string) {
	// Red + Panic + Exit(1)
	l.logWithColor(Red, "PANIC", msg)
	panic(msg)
}

// Formatted Trace log. Use this like fmt.Printf
func (l *LoggerImpl) Tracef(format string, args ...interface{}) {
	l.logWithColorf(Cyan, "TRACE", format, args...)
}

// Formatted Debug log. Use this like fmt.Printf
func (l *LoggerImpl) Debugf(format string, args ...interface{}) {
	// Check if runtime is development
	if os.Getenv("RUNTIME") != "development" {
		return
	}

	// Blue
	l.logWithColorf(Blue, "DEBUG", format, args...)
}

// Formatted Log log. Use this like fmt.Printf
func (l *LoggerImpl) Logf(format string, args ...interface{}) {
	// Green
	l.logWithColorf(Green, "LOG", format, args...)
}

// Formatted Warn log. Use this like fmt.Printf
func (l *LoggerImpl) Warnf(format string, args ...interface{}) {
	// Yellow
	l.logWithColorf(Yellow, "WARN", format, args...)
}

// Formatted Error log. Use this like fmt.Printf
func (l *LoggerImpl) Errorf(format string, args ...interface{}) {
	// Red
	l.logWithColorf(Red, "ERROR", format, args...)
}

// Formatted Fatal log. Use this like fmt.Printf
func (l *LoggerImpl) Fatalf(format string, args ...interface{}) {
	// Red + Fatal + Exit(1)
	l.logWithColorf(Red, "FATAL", format, args...)
	log.Fatalf(format, args...)
}

// Formatted Panic log. Use this like fmt.Printf
func (l *LoggerImpl) Panicf(format string, args ...interface{}) {
	// Red + Panic + Exit(1)
	l.logWithColorf(Red, "PANIC", format, args...)
	panic(fmt.Sprintf(format, args...))
}
