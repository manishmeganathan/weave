package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// A function to create a new directory given the path to the directory.
// If the directory already exists, this is a no-op.
func CreateDirectory(dirpath string) {
	// Check if the directory exists
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		// Create the directory
		err = os.MkdirAll(dirpath, 0755)
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to create directory: %v.", dirpath)
		}
	}
}

// A function to clear the contents of a directory given the path to the directory.
// If the directory does not exist, this is a no-op.
func ClearDirectory(dirpath string) {
	// Remove all contents of the directory
	err := os.RemoveAll(dirpath)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to clear directory.")
	}
}

// A function that returns the path to the config directory.
// The config directory is at %HOME%/blockweave/
func ConfigDirectory() string {
	// Create the path to the blockweave directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to detect home directory.")
	}

	// Return the path to the blockweave directory
	return filepath.Join(homedir, "blockweave")
}
