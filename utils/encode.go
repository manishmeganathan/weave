/*
This module contains the utility function that are used to
encode binary objects to different formats and to convert
structs into their binary gob formats
*/
package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"strconv"

	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
)

// A type alias for a byte slice that represents a gob of data
type Gob []byte

// An interface that defines a gob encodable object
// Must be serializable and deserializable as a Gob.
type GobEncodable interface {
	// A method that serializes the object into a gob
	Serialize() Gob
	// A method that deserializes the object from a gob
	Deserialize(Gob)
}

// A function to encode an object of arbirary type into a gob of bytes
func GobEncode(object interface{}) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer
	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the object into a gob
	err := encoder.Encode(object)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalf("failed to encode object of type %T as Gob.\n", object)
	}

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to decode a gob of bytes into an object of given type.
// The data of the object will be overriden with the gob data.
func GobDecode(gobdata Gob, object interface{}) interface{} {
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(object)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalf("failed to decode Gob as object of type %T.\n", object)
	}

	// Return the decoded object
	return object
}

// A function to encode a bytes payload into a Base58 bytes payload
func Base58Encode(payload []byte) []byte {
	// Encode the payload to base58
	encode := base58.Encode(payload)
	// Cast the encoded string to a slice of bytes and return it
	return []byte(encode)
}

// A function to decode a Base58 bytes payload into a regular bytes payload
func Base58Decode(encodeddata []byte) []byte {
	// Cast the base 58 encoded data into a string and decode it
	decode, err := base58.Decode(string(encodeddata[:]))
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Errorln("failed to decode from base58.")
	}

	// Return the decoded bytes
	return decode
}

// A function that encodes and returns an int as its hex/byte representation
func HexEncode(number int) []byte {
	// Format the integer into a string
	strint := strconv.FormatInt(int64(number), 10)
	// Convert the string into a slice of bytes
	src := []byte(strint)
	// Create a null destination object with
	// capacity for the encoded object
	dst := make([]byte, hex.EncodedLen(len(src)))

	// Encode the number into a hex
	hex.Encode(dst, src)
	// Return the hex value
	return dst
}

// A function that decodes and returns an int as from its hex/byte representation
func HexDecode(src []byte) int {
	// Create a null destination object with
	// capacity for the decoded object
	dst := make([]byte, hex.DecodedLen(len(src)))

	// Decode the number from a hex
	_, err := hex.Decode(dst, src)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Errorln("failed to decode from hexadecimal.")
	}

	// Convert the decoded integer bytes into a string
	strint := string(dst)
	// Parse the string into an integer
	number, err := strconv.ParseInt(strint, 10, 0)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Errorln("failed to parse integer from hexadecimal.")
	}

	// Return the integer
	return int(number)
}
