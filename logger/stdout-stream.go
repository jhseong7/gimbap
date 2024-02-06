package logger

import (
	"fmt"
)

type (
	StdOutStream struct {
		ILogStream
	}
)

func (s *StdOutStream) Write(msg LogMessage) {
	// Print the log message to the console
	fmt.Printf(
		"%s[%s] %s %6s [%s] %s%s\n", // Format string
		msg.Color,                   // Set colour
		msg.AppName,                 // Add the Prefix (app name)
		msg.Time,                    // Add the time
		msg.Level,                   // Add the log level
		msg.Name,                    // Add the log name
		msg.Msg,                     // Add the message
		Reset,                       // Reset colour
	)
}

func NewStdOutStream() *StdOutStream {
	return &StdOutStream{}
}
