package blockchain

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
)

const dbfile = "./db/blocks/MANIFEST"

// A function that generates and returns the
// Hex/Bytes representation of an int64
func Hexify(number int64) []byte {
	// Construct a new binary buffer
	buff := new(bytes.Buffer)
	// Write the number as a binary into the buffer in Big Endian order
	err := binary.Write(buff, binary.BigEndian, number)
	// Handle any potential error
	Handle(err)

	// Return the bytes from the binary buffer
	return buff.Bytes()
}

// A function to handle errors.
func Handle(err error) {
	// Check if error is non null
	if err != nil {
		// Log the error and throw a panic
		log.Panic(err)
	}
}

// A function to check if the DB exists
func DBexists() bool {
	// Check if the database MANIFEST file exists
	if _, err := os.Stat(dbfile); errors.Is(err, os.ErrNotExist) {
		// Return false if the file is missing
		return false
	}

	// Return true if the file exists
	return true
}
