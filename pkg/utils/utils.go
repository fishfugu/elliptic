package utils

import (
	"math/big"
	"os"

	"github.com/sirupsen/logrus"
)

const outputDebugInfo = true

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

// FindPrime finds a prime number within the range specified by the strings i and high.
// Returns the first prime found as a string, or an error if no prime is found or inputs are invalid.
func FindPrime(i, high *big.Int) *big.Int {
	// TODO: is there any reason to implement different way / direction of search?
	// Start searching for a prime at the low end of the range
	for p := i; new(big.Int).Set(p).Cmp(high) < 0; p = new(big.Int).Add(p, new(big.Int).SetInt64(1)) {
		if p.ProbablyPrime(1000) {
			return p
		}
	}
	return new(big.Int).SetInt64(0)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// newFloat creates a new big.Float with the default precision
// ue this whenever creating new Float (except in very specific circumstances...)
func NewFloat() *big.Float {
	return new(big.Float).SetPrec(0)
}
