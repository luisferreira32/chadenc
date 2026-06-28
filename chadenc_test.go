package chadenc

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func Test_ChandEncode(t *testing.T) {
	testcases := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  error
	}{
		{
			name:     "simple case",
			input:    []byte{1, 2, 3},
			expected: []byte{1, 1, 2, 2, 3, 3},
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "single byte",
			input:    []byte{42},
			expected: []byte{42, 42},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := &Encoder{}
			output, err := encoder.EncodeFrame(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error: %v, want: %v", err, tc.wantErr)
			}
			if string(output) != string(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, output)
			}
		})
	}
}

func Test_ChandEncodeStream(t *testing.T) {
	testcases := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  error
	}{
		{
			name:     "simple case",
			input:    []byte{1, 2, 3},
			expected: []byte{1, 1, 2, 2, 3, 3},
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "single byte",
			input:    []byte{42},
			expected: []byte{42, 42},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := &Encoder{}
			output := &bytes.Buffer{}
			input := &bytes.Buffer{}
			input.Write(tc.input)

			err := encoder.EncodeStream(input, output)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error: %v, want: %v", err, tc.wantErr)
			}
			if output.String() != string(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, output.String())
			}
		})
	}
}

func Test_ChandDecode(t *testing.T) {
	testcases := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  error
	}{
		{
			name:     "simple case",
			input:    []byte{1, 1, 2, 2, 3, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "single byte",
			input:    []byte{42, 42},
			expected: []byte{42},
		},
		{
			name:    "odd length input",
			input:   []byte{1, 1, 2},
			wantErr: io.ErrUnexpectedEOF,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := &Decoder{}
			output, err := decoder.DecodeFrame(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error: %v, want: %v", err, tc.wantErr)
			}
			if string(output) != string(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, output)
			}
		})
	}
}

func Test_ChandDecodeStream(t *testing.T) {
	testcases := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  error
	}{
		{
			name:     "simple case",
			input:    []byte{1, 1, 2, 2, 3, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "single byte",
			input:    []byte{42, 42},
			expected: []byte{42},
		},
		{
			name:     "multiple frames",
			input:    []byte{1, 1, 2, 2, 3, 3, 4, 4},
			expected: []byte{1, 2, 3, 4},
		},
		{
			name:     "odd length input",
			input:    []byte{1, 1, 2},
			expected: []byte{1},
			wantErr:  io.ErrUnexpectedEOF,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := &Decoder{BufferSize: 2} // small buffer size to test multiple reads
			output := &bytes.Buffer{}
			input := &bytes.Buffer{}
			input.Write(tc.input)

			err := decoder.DecodeStream(input, output)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error: %v, want: %v", err, tc.wantErr)
			}
			if output.String() != string(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, output.String())
			}
		})
	}
}
