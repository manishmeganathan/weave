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
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/sirupsen/logrus"
	"github.com/vrecan/death/v3"
)

// A function to check if the DB exists
func DBExists() bool {
	// Check if the database MANIFEST file exists
	if _, err := os.Stat(utils.DBfile); errors.Is(err, os.ErrNotExist) {
		// Return false if the file is missing
		return false
	}
	// Return true if the file exists
	return true
}

type DatabaseClient struct {
	Client *badger.DB
	IsOpen bool
}

// A constructor function that generates and returns
// a new Database object that has been opened
func NewDatabaseClient() *DatabaseClient {
	// Define the BadgerDB options for the DB path
	opts := badger.DefaultOptions(utils.DBpath)
	// Switch off the Badger Logger
	opts.Logger = nil

	// Construct an emtpy Database object
	db := &DatabaseClient{Client: nil, IsOpen: false}
	// Open the database
	db.Open(opts)

	// Setup database to close at application death
	go db.CloseOnDeath()

	// Return the database
	return db
}

// A method of Database that opens the BadgerDB
// client with the given badger DB options
func (db *DatabaseClient) Open(opts badger.Options) {
	// Open the Badger DB with the defined options
	client, err := badger.Open(opts)
	// Handle any potential error
	utils.HandleErrorLog(err, "database client failed to open")

	// Assign the DB client
	db.Client = client
	// Set the open flag to true
	db.IsOpen = true
	// log the opening of the database
	logrus.Info("database client has been opened")
}

// A method of Database that closes the BadgerDB client
func (db *DatabaseClient) Close() {
	// log the closing of the database
	logrus.Info("database client has been closed")
	// Empty the client field
	db.Client = nil
	// Set the open flag to false
	db.IsOpen = false

	// Close the client
	db.Client.Close()
}

// A method of Database that closes the connection of the
// BadgerDB client upon the runtime closing abruptly
func (db *DatabaseClient) CloseOnDeath() {
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

// A method of Database that deletes all entries
// with a given prefix from the Badger DB.
func (db *DatabaseClient) DeleteKeyPrefix(prefix []byte) {

	// Define a function that accepts a 2D slice of byte keys to delete
	DeleteKeys := func(keystodelete [][]byte) error {

		// Define an Update transaction on the database
		err := db.Client.Update(func(txn *badger.Txn) error {
			// Iterate over the keys to delete
			for _, key := range keystodelete {
				// Delete the key
				if err := txn.Delete(key); err != nil {
					// Return any potential error
					return err
				}
			}
			// Return nil error
			return nil
		})

		// Check if key deletion transaction has comlpeted without error
		if err != nil {
			// Return the error
			return err
		}
		// Return nil error
		return nil
	}

	// Define the size limit of key accumulation. This value
	// is based on BadgerDBs optimal object handling limit.
	collectlimit := 100000

	// Define a View transaction on the database
	db.Client.View(func(txn *badger.Txn) error {

		// Set up the default iteration options for the database
		opts := badger.DefaultIteratorOptions
		// Set value pre-fetching to off (We only need check the key's prefix to accumulate it)
		opts.PrefetchValues = false
		// Create an iterator with the options
		dbiterator := txn.NewIterator(opts)
		// Defer the closing of the iterator
		defer dbiterator.Close()

		// Declare collection counter
		keyscollected := 0
		// Create a 2D slice of bytes to collect keys
		keystodelete := make([][]byte, 0)

		// Start the iterator and seek the keys with the provided prefix (validate the keys for the prefix as well)
		for dbiterator.Seek(prefix); dbiterator.ValidForPrefix(prefix); dbiterator.Next() {

			// Make a copy of the key value
			key := dbiterator.Item().KeyCopy(nil)
			// Add the key to slice
			keystodelete = append(keystodelete, key)
			// Increment the counter
			keyscollected++

			// Check if counter has reached the collection limit
			if keyscollected == collectlimit {
				// Delete all keys accumulated so far
				if err := DeleteKeys(keystodelete); err != nil {
					// Handle any potential error
					logrus.WithFields(logrus.Fields{
						"prefix": prefix,
						"error":  err,
					}).Fatal("database key prefix deletion failed!")
				}

				// Reset the key accumulation
				keystodelete = make([][]byte, 0)
				// Reset the key accumulation counter
				keyscollected = 0
			}
		}

		// Check if there any keys that have been collected (but less than the collection limit)
		if keyscollected > 0 {
			// Delete all the accumulated keys
			if err := DeleteKeys(keystodelete); err != nil {
				// Handle any potential errors
				logrus.WithFields(logrus.Fields{
					"prefix": prefix,
					"error":  err,
				}).Fatal("database key prefix deletion failed!")
			}
		}

		// Return a nil error for the transaction
		return nil
	})
}
