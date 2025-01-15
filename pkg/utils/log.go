package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitializeLogger sets up the logger and returns it for use in other parts of the application.
func InitialiseLogger(prefix string) *logrus.Logger {
	// Create a logger that writes to stdout
	var log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.WarnLevel,
	}
	return log
}

// LogOnError logs an error consistently with a given context message.
func LogOnError(logger *logrus.Logger, err error, context string, panicOnError bool) {
	if err != nil {
		logger.Printf("[ERROR] %s: %v", context, err)
		if panicOnError {
			logrus.Panicf("[ERROR] %s: %v", context, err)
		}
	}
}

// WarnOnError logs an error consistently with a given context message.
func WarnOnError(logger *logrus.Logger, err error, context string) {
	if err != nil {
		logger.Printf("[WARN] %s: %v", context, err)
	}
}

// LogFailure logs a failure with a provided boolean and context message.
func LogOnFailure(logger *logrus.Logger, success bool, context string) {
	if !success {
		logger.Printf("[FAILURE] %s", context)
	}
}
