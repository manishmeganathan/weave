package utils

import (
	"bytes"
	"testing"
)

func Test_Base58Encode(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello", "Cn8eVZg"},
		{"world", "EUYUqQf"},
		{"1lL 0oO", "2sdUXQxyyx"},
	}

	for _, tt := range tests {
		data := []byte(tt.input)
		encoded := Base58Encode(data)

		if string(encoded) != tt.output {
			t.Fatalf("Base58Encode('%v') failed! expected: %v, got: %v", tt.input, tt.output, string(encoded))
		}
	}
}

func Test_Base58Decode(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"Cn8eVZg", "hello"},
		{"EUYUqQf", "world"},
		{"2sdUXQxyyx", "1lL 0oO"},
	}

	for _, tt := range tests {
		data := []byte(tt.input)
		decoded := Base58Decode(data)

		if string(decoded) != tt.output {
			t.Fatalf("Base58Decode('%v') failed! expected: %v, got: %v", tt.input, tt.output, string(decoded))
		}
	}
}

func Test_HexEncode(t *testing.T) {
	tests := []struct {
		input  int
		output []byte
	}{
		{0, []byte{51, 48}},
		{1, []byte{51, 49}},
		{100, []byte{51, 49, 51, 48, 51, 48}},
		{4000, []byte{51, 52, 51, 48, 51, 48, 51, 48}},
		{52235, []byte{51, 53, 51, 50, 51, 50, 51, 51, 51, 53}},
	}

	for _, tt := range tests {
		encoded := HexEncode(tt.input)

		if !bytes.Equal(encoded, tt.output) {
			t.Fatalf("HexEncode(%v) failed! expected: %v, got: %v", tt.input, tt.output, encoded)
		}
	}
}

func Test_HexDecode(t *testing.T) {
	tests := []struct {
		input  []byte
		output int
	}{
		{[]byte{51, 48}, 0},
		{[]byte{51, 49}, 1},
		{[]byte{51, 49, 51, 48, 51, 48}, 100},
		{[]byte{51, 52, 51, 48, 51, 48, 51, 48}, 4000},
		{[]byte{51, 53, 51, 50, 51, 50, 51, 51, 51, 53}, 52235},
	}

	for _, tt := range tests {
		decoded := HexDecode(tt.input)

		if decoded != tt.output {
			t.Fatalf("HexDecode(%v) failed! expected: %v, got: %v", tt.input, tt.output, decoded)
		}
	}
}
