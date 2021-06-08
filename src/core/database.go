/*
This module contains the definition and implementation
of the Database structure and its methods
*/
package core

import (
	"errors"
	"os"
	"runtime"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/sirupsen/logrus"
	"github.com/vrecan/death/v3"
)

type Database struct {
	Client *badger.DB
	IsOpen bool
}

// TODO: MOVE TO A LOCATION OUTSIDE OF THE REPO DIR
const dbfile string = "./tmp/db/blocks/MANIFEST"
const dbpath string = "./tmp/db/blocks"

// A function to check if the DB exists
func DBExists() bool {
	// Check if the database MANIFEST file exists
	if _, err := os.Stat(dbfile); errors.Is(err, os.ErrNotExist) {
		// Return false if the file is missing
		return false
	}
	// Return true if the file exists
	return true
}

// A method of Database that opens the BadgerDB
// client with the given badger DB options
func (db *Database) Open(opts badger.Options) {
	// Open the Badger DB with the defined options
	client, err := badger.Open(opts)
	// Handle any potential error
	logrus.Fatal("database client failed to open", err)

	// Assign the DB client
	db.Client = client
	// Set the open flag to true
	db.IsOpen = true
}

// A method of Database that closes the BadgerDB client
func (db *Database) Close() {
	// Set the open flag to false
	db.IsOpen = false
	// Close the client
	db.Client.Close()
}

// A method of Database that closes the connection of the
// BadgerDB client upon the runtime closing abruptly
func (db *Database) CloseOnDeath() {
	// Setup death signals
	demise := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Anon function that executes when a death signal is triggered
	demise.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()

		// Log the death of the database client
		logrus.Infoln("connection to database has terminated due to application death")

		// Close the database client
		db.Close()
	})
}

// A constructor function that generates and returns
// a new Database object that has been opened
func NewDatabase() *Database {
	// Define the BadgerDB options for the DB path
	opts := badger.DefaultOptions(dbpath)
	// Switch off the Badger Logger
	opts.Logger = nil

	// Construct an emtpy Database object
	db := &Database{Client: nil, IsOpen: false}
	// Open the database
	db.Open(opts)

	// Setup database to close at application death
	go db.CloseOnDeath()

	// Return the database
	return db
}
