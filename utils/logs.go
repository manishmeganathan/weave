package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// A function that initializes the log level, formatter and output
func LogInitialize(loglevel logrus.Level) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the sepcified level or above.
	logrus.SetLevel(loglevel)
}
