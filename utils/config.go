package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Represents the version of the source code
const SrcVersion = "0.5.0"

var (
	// Represents the prefix key used for utxo keys
	UTXOprefix = []byte("utxo-")
	// Represents the key used for storing chain head
	ChainHeadKey = []byte("chainhead")
	// Represents the key used for storing the chain height
	ChainHeightKey = []byte("chainheight")
)

// A struct that represents the contents of the config file.
type Config struct {
	// Represents the jbok configuration
	JBOK jbokconfig `json:"jbok"`
	// Represents the database configuration
	DB dbconfig `json:"db"`
}

// A struct that represents a jbok configuration
type jbokconfig struct {
	// Represents the path to the JBOK data file
	File string `json:"file"`
	// Represents the default address used for mining rewards
	Default string `json:"default"`
}

// A struct that represents a database configuration
type dbconfig struct {
	// Represents the path to the database directory
	Root string `json:"root"`
	// Represents the configuration of the State bucket
	State bucketconfig `json:"state"`
	// Represents the configuration of the Blocks bucket
	Blocks bucketconfig `json:"blocks"`
}

// A struct that represents a database bucket configuration
type bucketconfig struct {
	// Represents the path to the bucket manifest file
	File string `json:"file"`
	// Represents the path to the bucket directory
	Directory string `json:"directory"`
}

// A function that returns the path to the config file.
// The config file is at %HOME%/blockweave/config.json
func getconfigfilepath() string {
	// Retrieve the path to the config dir.
	configdir := ConfigDirectory()
	// Return the file location
	return filepath.Join(configdir, "config.json")
}

// A function that checks if the config file exists in the expected location.
// Returns an error if the file does not exist.
func CheckConfigFile() error {
	// Get the path to the config file.
	filelocation := getconfigfilepath()

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
func ReadConfigFile() *Config {
	// Get the path to the config file.
	filelocation := getconfigfilepath()

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
	return &config
}

// A function that generates a new config file with the default
// values, writes it to the file and returns it as a Config object.
//
// If the file already exists, it will be overwritten.
// If the blockweave directory does not exist, it will be created.
func GenerateConfigFile(write bool) *Config {
	// Get the path to the blockweave directory.
	configdir := ConfigDirectory()
	// Create the blockweage directory if it does not exist.
	CreateDirectory(configdir)

	// Generate a default Config with default values.
	defaultconfig := Config{
		JBOK: jbokconfig{
			File:    filepath.Join(configdir, "jbok.data"),
			Default: "",
		},
		DB: dbconfig{
			Root: filepath.Join(configdir, "db"),
			State: bucketconfig{
				File:      filepath.Join(configdir, "db", "state", "MANIFEST"),
				Directory: filepath.Join(configdir, "db", "state"),
			},
			Blocks: bucketconfig{
				File:      filepath.Join(configdir, "db", "blocks", "MANIFEST"),
				Directory: filepath.Join(configdir, "db", "blocks"),
			},
		},
	}

	// Check if write flag is set
	if write {
		// Write the generated config and check for errors.
		err := defaultconfig.WriteConfigFile()
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to write generated config file.")
		}
	}

	// Return the default config
	return &defaultconfig
}

// A method of Config that write the config data to the file.
// If the file already exists, it will be overwritten.
func (config *Config) WriteConfigFile() error {
	// Get the path to the config file.
	filelocation := getconfigfilepath()

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

// A method of Config that prints the configuration
// values in a formatted menu-list to stdout.
func (config *Config) PrintConfigFile() {
	// Print out the
	fmt.Println()
	fmt.Println("-----Weave-Configuration-File-----")
	fmt.Println()

	fmt.Println("----JBOK-Configuration----")
	fmt.Printf("JBOK File: %v\n", config.JBOK.File)
	fmt.Printf("JBOK Default: %v\n", config.JBOK.Default)
	fmt.Println()

	fmt.Println("----Database-Configuration----")
	fmt.Printf("DB Root: %v\n", config.DB.Root)
	fmt.Printf("DB State File: %v\n", config.DB.State.File)
	fmt.Printf("DB State Directory: %v\n", config.DB.State.Directory)
	fmt.Printf("DB Blocks File: %v\n", config.DB.Blocks.File)
	fmt.Printf("DB Blocks Directory: %v\n", config.DB.Blocks.Directory)
	fmt.Println()

	// fmt.Println("----Network-Configuration----")
	// fmt.Println()

	fmt.Println("----end-of-file----")
	fmt.Println()
}
