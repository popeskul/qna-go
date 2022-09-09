// Package logger defines the logger for the application.
// It's a singleton.
// It's used to log messages to the console and file.
package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

// Logger an entry is the final or intermediate Logrus logging entry.
type Logger struct {
	*logrus.Entry
}

// writerHook is a hook that writes logs to a writer.
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire is the function that is called when a log is fired.
// It writes the log to the file and the console.
// It returns the error if any.
func (h *writerHook) Fire(entry *logrus.Entry) error {
	s, err := entry.String()
	if err != nil {
		return err
	}

	for _, writer := range h.Writer {
		_, err = writer.Write([]byte(s))
		if err != nil {
			return err
		}
	}

	return nil
}

// Levels returns the available logging levels.
func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}

// GetLogger creates a new instance of logger if it doesn't exist
func GetLogger() *Logger {
	return &Logger{e}
}

// Init initializes the logger.
// It sets the log level and the output.
// It returns the logger.
// It panics if the log level is invalid.
// It panics if the output is invalid.
// It panics if the file is invalid.
func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println("Getwd->", path)

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	allLogsFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allLogsFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
