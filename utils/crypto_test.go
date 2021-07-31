package utils

import (
	"bytes"
	"testing"
)

func Test_Hash256(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte("hello"),
			[]byte{13, 144, 59, 232, 5, 20, 39, 4, 242, 18, 133, 211, 4, 206, 74, 177, 32, 173, 137, 99, 229, 152, 161, 136, 49, 67, 27, 168, 91, 215, 246, 61},
		},
		{
			[]byte("world"),
			[]byte{114, 118, 135, 42, 94, 245, 120, 188, 166, 208, 180, 114, 225, 0, 153, 204, 87, 33, 120, 75, 168, 228, 246, 174, 9, 174, 128, 248, 166, 36, 196, 224},
		},
	}

	for _, tt := range tests {
		hash := Hash256(tt.input)

		if len(hash) != 32 {
			t.Fatalf("incorrect hash length! expected: 32, got: %v", len(hash))
		}

		if !bytes.Equal(hash, tt.output) {
			t.Fatalf("incorrect hash output! expected: %v, got: %v", tt.output, hash)
		}
	}
}

func Test_Hash160(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte("hello"),
			[]byte{197, 45, 112, 90, 72, 192, 84, 218, 125, 241, 223, 182, 49, 52, 19, 160, 240, 23, 89, 159},
		},
		{
			[]byte("world"),
			[]byte{208, 169, 89, 196, 244, 220, 252, 127, 188, 74, 206, 158, 64, 24, 94, 101, 107, 173, 226, 69},
		},
	}

	for _, tt := range tests {
		hash := Hash160(tt.input)

		if len(hash) != 20 {
			t.Fatalf("incorrect hash length! expected: 20, got: %v", len(hash))
		}

		if !bytes.Equal(hash, tt.output) {
			t.Fatalf("incorrect hash output! expected: %v, got: %v", tt.output, hash)
		}
	}
}

func Test_Hash32(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		{
			[]byte("hello"),
			[]byte{13, 144, 59, 232},
		},
		{
			[]byte("world"),
			[]byte{114, 118, 135, 42},
		},
	}

	for _, tt := range tests {
		hash := Hash32(tt.input)

		if len(hash) != 4 {
			t.Fatalf("incorrect hash length! expected: 4, got: %v", len(hash))
		}

		if !bytes.Equal(hash, tt.output) {
			t.Fatalf("incorrect hash output! expected: %v, got: %v", tt.output, hash)
		}
	}
}
