package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

// Represents the version of the source code
const SrcVersion = "0.5.0"

// Represents the prefix key used for utxo keys
var UTXOprefix = []byte("utxo-")

// A struct that represents a config file and its data
type Config struct {
	// Represents the path to the JBOK wallet data file
	JBOKFile string `json:"jbokfile"`
	// Represents the path to the BadgerDB manifest file
	DBFile string `json:"dbfile"`
	// Represents the path to the chain database directory
	DBDirectory string `json:"dbdir"`
}

// A function that returns the path to the config file.
// The config file is stored at %HOME%/blockweave/config.json
func getconfigpath() string {
	// Retrieve the path to the config file.
	filelocation, err := homedir.Expand("~/blockweave/config.json")
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to detect home directory.")
	}

	// Return the file location
	return filelocation
}

// A function that checks if the config file exists in the expected location.
// Returns an error if the file does not exist.
func CheckConfigFile() error {
	// Get the path to the config file.
	filelocation := getconfigpath()

	// Check if the file exists at the location
	if _, err := os.Stat(filelocation); err == nil {
		// File exists.
		return nil
	} else if os.IsNotExist(err) {
		// File does not exist.
		return fmt.Errorf("config file does not exist")
	} else {
		// File may or may not exist.
		return fmt.Errorf("could not determine if config file exists")
	}
}

// A function that reads the config file and returns the data as a Config object.
func ReadConfigFile() Config {
	// Get the path to the config file.
	filelocation := getconfigpath()

	// Open the config file
	configfile, err := os.Open(filelocation)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to open config file.")
	}

	// Defer the closing of the file
	defer configfile.Close()

	// Read the config file into a byte array
	var config Config
	byteValue, _ := ioutil.ReadAll(configfile)

	// Unmarshal the JSON byte array into a struct and return it
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

// A function that writes a given Config object to the config file.
// If the file already exists, it will be overwritten.
func WriteConfigFile(config Config) error {
	// Get the path to the config file.
	filelocation := getconfigpath()

	// Format and indent the config object provided into a byte array.
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("could not format and marshal config. %v", err)
	}

	// Write the byte array to the file location.
	if err = ioutil.WriteFile(filelocation, file, 0644); err != nil {
		return fmt.Errorf("could not write config. %v", err)
	}

	return nil
}

// A function that generates a new config file with the default values.
// If the file already exists, it will be overwritten.
// If the blockweave directory does not exist, it will be created.
func GenerateConfigFile() {
	// Create the path to the blockweave directory
	configdir, err := homedir.Expand("~/blockweave")
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to detect home directory.")
	}

	// Check if the blockweave directory exists
	_, err = os.Stat(configdir)
	if os.IsNotExist(err) {
		// Create the blockweave directory
		err = os.MkdirAll(configdir, 0755)
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to create blockweave directory.")
		}
	}

	// Generate a default Config with default values.
	defaultConfig := Config{
		JBOKFile:    configdir + "/jbok.data",
		DBFile:      configdir + "/db/manifest.json",
		DBDirectory: configdir + "/db/",
	}

	// Write the generated config and check for errors.
	err = WriteConfigFile(defaultConfig)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to write generated config file.")
	}
}
