package persistence

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/sirupsen/logrus"
	"github.com/vrecan/death/v3"
)

// A type alias that represents a type of database bucket
type Bucket uint32

// A set of constants that represent valids types of database buckets
const (
	STATE Bucket = iota
	BLOCKS
)

// A struct that represents the client for a database bucket
type DatabaseBucket struct {
	// Represents the BadgerDB client for the bucket
	Client *badger.DB
	// Represents the type of bucket
	Bucket Bucket
	// Represents the whether the client is open
	IsOpen bool
}

// A function to check if the database exists locally.
// Checks if all database buckets exist by confirming
// the existence of the MANIFEST file for each bucket.
//
// If either MANIFEST file does not exist, the function
// clears the contents of the database root and creates
// the directories for the database buckets.
func CheckDatabase() bool {
	// Get the Config data
	config := utils.ReadConfigFile()

	// Retrieve the file status for the database manifest files for the buckets
	_, err_state := os.Stat(config.DB.State.File)
	_, err_blocks := os.Stat(config.DB.Blocks.File)

	// Check if either bucket does not exist
	if errors.Is(err_state, os.ErrNotExist) || errors.Is(err_blocks, os.ErrNotExist) {
		// Clear the contents of the db root directory
		utils.ClearDirectory(config.DB.Root)

		// Create an empty state db directory if it does not exist
		utils.CreateDirectory(config.DB.State.Directory)
		// Create an empty blocks db directory if it does not exist
		utils.CreateDirectory(config.DB.Blocks.Directory)

		// Return false because some file does not exist
		return false
	}

	// Return true because all bucket files exist
	return true
}

// A constructor function that generates and returns
// a new Database bucket object that has been opened
// The bucket argument is the type of bucket to open
// Valid options are the STATE and BLOCKS constants.
func NewDatabaseBucket(bucket Bucket) *DatabaseBucket {
	// Get the Config data
	config := utils.ReadConfigFile()

	// Declare a new badger options variable
	var opts badger.Options
	// Check the type of bucket
	switch bucket {
	// The chain state bucket
	case STATE:
		// Set the Badger DB options for the state bucket
		opts = badger.DefaultOptions(config.DB.State.Directory)

	// The chain blocks bucket
	case BLOCKS:
		// Set the Badger DB options for the blocks bucket
		opts = badger.DefaultOptions(config.DB.Blocks.Directory)

	// Invalid type
	default:
		// Log the fatal error
		logrus.WithFields(logrus.Fields{"error": errors.New("invalid bucket type")}).Fatalln("failed to create database bucket client.")
	}

	// Switch off the Badger Logger
	opts.Logger = nil

	// Construct an empty database bucket object
	db := &DatabaseBucket{Client: nil, Bucket: bucket, IsOpen: false}
	// Open the database
	db.Open(opts)

	// Setup database to close at application death
	go db.safedeath()

	// Return the database
	return db
}

// A method of DatabaseBucket that opens the BadgerDB
// client for the db bucket with the given badger DB options
func (db *DatabaseBucket) Open(opts badger.Options) {
	// Open the Badger DB bucket with the defined options
	client, err := badger.Open(opts)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to open database bucket.")
	}

	// Assign the DB client
	db.Client = client
	// Set the open flag to true
	db.IsOpen = true

	// log the opening of the database
	logrus.Infof("database %v bucket client has been opened\n", db.Bucket)
}

// A method of DatabaseBucket that closes the BadgerDB client for the bucket
func (db *DatabaseBucket) Close() {
	// log the closing of the database
	logrus.Infof("database %v bucket client has been closed\n", db.Bucket)
	// Close the client
	db.Client.Close()

	// Empty the client field
	db.Client = nil
	// Set the open flag to false
	db.IsOpen = false
}

// A method of DatabaseBucket that closes the connection of the
// BadgerDB client for the bucket when the runtime closes abruptly
func (db *DatabaseBucket) safedeath() {
	// Setup death signals
	demise := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Lambda function that executes when a death signal is triggered
	demise.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()

		// Log the death of the database client
		logrus.Infoln("connection to database has terminated due to application death")

		// Close the database client
		db.Close()
	})
}

// A method of DatabaseBucket that retrieves the value for
// a given key from the BadgerDB client for the bucket
func (db *DatabaseBucket) GetKey(key []byte) ([]byte, error) {
	// Declare a variable to hold the value
	var value []byte

	// Define a view transaction on the database bucket
	err := db.Client.View(func(txn *badger.Txn) error {
		// Get the item with the key from the DB
		item, err := txn.Get(key)
		// Return any potential error
		if err != nil {
			return fmt.Errorf("failed to GET database item! error - %v", err)
		}

		// Retrieve the value of the item
		if err = item.Value(func(val []byte) error {
			// Set the value
			value = val
			return nil

			// Return any potential error
		}); err != nil {
			return fmt.Errorf("failed to GET database value! error - %v", err)
		}

		// Return the nil error
		return nil
	})

	// Return the value and any error generated
	return value, err
}

// A method of DatabaseBucket that sets the value for a given key-value pair
func (db *DatabaseBucket) SetKey(key, value []byte) error {
	// Define an update transaction on the database bucket
	err := db.Client.Update(func(txn *badger.Txn) error {
		// Add the key-value pair to the database
		if err := txn.Set(key, value); err != nil {
			// Return any potential error
			return fmt.Errorf("failed to SET database key! error - %v", err)
		}

		// Return the nil error
		return nil
	})

	// Return any error generated
	return err
}

// A method of DatabaseBucket that deletes all entries
// with a given prefix from the Badger DB bucket.
func (db *DatabaseBucket) DeleteKeyPrefix(prefix []byte) {

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
					}).Fatal("failed to delete database key prefix!")
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
				}).Fatal("failed to delete database key prefix!")
			}
		}

		// Return a nil error for the transaction
		return nil
	})
}
