package logger

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type (
	FileLogStreamOption struct {
		LogDirectory  string
		FileName      string
		FileRollover  bool
		MaxFileSizeKb int
	}

	FileLogStream struct {
		ILogStream

		// Copy of the initial options
		options FileLogStreamOption

		// File pointer
		file *os.File

		// Mutex to prevent multiple writes at the same time
		mutex *sync.Mutex
	}
)

func (s *FileLogStream) getLogFileName(prefix, date string) string {
	if s.options.FileRollover {
		return fmt.Sprintf("%s.%s.log", prefix, date)
	}

	return fmt.Sprintf("%s.log", prefix)
}

// NOTE: Use this later when rollover is implemented
// func (s *FileLogStream) getFileSizeKb() int64 {
// 	fileInfo, err := s.file.Stat()

// 	if err != nil {
// 		panic(err)
// 	}

// 	return fileInfo.Size() / 1024
// }

func (s *FileLogStream) createNewLog() {
	// Close the file
	s.file.Close()

	// If the LogDirectory does not exist, create it (recursively)
	if _, err := os.Stat(s.options.LogDirectory); os.IsNotExist(err) {
		os.MkdirAll(s.options.LogDirectory, 0755)
	}

	// Open a new file
	f, err := os.OpenFile(
		path.Join(s.options.LogDirectory, s.getLogFileName(s.options.FileName, time.Now().Format("2006-01-02"))),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		panic(err)
	}

	s.file = f
}

func (s *FileLogStream) getFilePointer() *os.File {
	// Get the current time, and if the current file is not the same date, close the file and open a new one with the current date
	currentDate := time.Now().Format("2006-01-02")
	logFileName := s.getLogFileName(s.options.FileName, currentDate)

	// If there is no file pointer, open a new file
	// If the file is not the same date, close the file and open a new one
	if s.file == nil || s.file.Name() != logFileName {
		s.createNewLog()
	}

	return s.file
}

func (s *FileLogStream) Write(msg LogMessage) {
	// Use a mutex to prevent multiple writes at the same time
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file := s.getFilePointer()

	// Append the message to the file (with a new line)
	// Format is "[prefix] time level [name] message"
	fs := fmt.Sprintf("[%s] %s %6s [%s] %s\n", msg.AppName, msg.Time, msg.Level, msg.Name, msg.Msg)
	fmt.Fprintln(file, fs)
}

func NewFileLogStream(option FileLogStreamOption) *FileLogStream {
	return &FileLogStream{
		options: option,
	}
}
