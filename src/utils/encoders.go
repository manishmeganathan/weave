/*
This module contains the utility function that are used to
encode binary objects to different formats and to convert
structs into their binary gob formats
*/
package utils

import (
	"bytes"
	"encoding/binary"

	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
)

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
	// Handle any potential errors
	logrus.Fatal("base58 decode failed!", err)
	// Return the decoded bytes
	return decode
}

// A function that encodes and returns an int64 as its hex/byte representation
func HexEncode(number int64) []byte {
	// Construct a new binary buffer
	buff := new(bytes.Buffer)
	// Write the number as a binary into the buffer in Big Endian order
	err := binary.Write(buff, binary.BigEndian, number)
	// Handle any potential error
	logrus.Fatal("integer hex encode failed!", err)

	// Return the bytes from the binary buffer
	return buff.Bytes()
}
