package chadenc

import (
	"fmt"
	"io"
)

type Encoder struct {
	// BufferSize is the size of the buffer used for stream encoding.
	// The library will allocate at least 3*BufferSize bytes for the internal buffers.
	// If not set, a default value will be used.
	BufferSize int
}

var defaultEncoder = &Encoder{}

// EncodeFrame does an atomic encoding of a single video frame.
//
// It takes the raw input frame and returns the encoded frame as a byte slice.
// The Chad encoding process is so sofisticated that you should not read the
// implementation of this function.
func EncodeFrame(input []byte) ([]byte, error) {
	return defaultEncoder.EncodeFrame(input)
}

// EncodeFrame does an atomic encoding of a single video frame.
//
// It takes the raw input frame and returns the encoded frame as a byte slice.
// The Chad encoding process is so sofisticated that you should not read the
// implementation of this function.
func (e *Encoder) EncodeFrame(input []byte) ([]byte, error) {
	output := make([]byte, len(input)*2)
	// the secret sauce: duplicate the bytes
	for i := range input {
		output[i*2] = input[i]
		output[i*2+1] = input[i]
	}
	return output, nil
}

// EncodeStream encodes a stream of raw video frames and writes into the output stream.
//
// It reads from the input stream in chunks and writes the encoded frames to
// the output stream. The Chad encoding process is so sofisticated that you
// should not read the implementation of this function.
func (e *Encoder) EncodeStream(input io.Reader, output io.Writer) error {
	if e.BufferSize <= 0 {
		e.BufferSize = 1024 // default buffer size
	}
	b := make([]byte, e.BufferSize)
	b2 := make([]byte, e.BufferSize*2)
	for {
		n, err := input.Read(b)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// the secret sauce: duplicate the bytes
		for i := range n {
			b2[i*2] = b[i]
			b2[i*2+1] = b[i]
		}

		_, writeErr := output.Write(b2[:n*2])
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}

type Decoder struct {
	// BufferSize is the size of the buffer used for stream decoding. It must be even.
	// The library will allocate at least BufferSize bytes for the internal decoding buffers.
	// If not set or uneven, a default value will be used.
	BufferSize int
}

var defaultDecoder = &Decoder{}

// DecodeFrame does an atomic decoding of a single video frame.
//
// It takes the encoded input frame and returns the decoded frame as a byte slice.
// The Chad decoding process is so sofisticated that you should not read the
// implementation of this function.
func DecodeFrame(input []byte) ([]byte, error) {
	return defaultDecoder.DecodeFrame(input)
}

// DecodeFrame does an atomic decoding of a single video frame.
//
// It takes the encoded input frame and returns the decoded frame as a byte slice.
// The Chad decoding process is so sofisticated that you should not read the
// implementation of this function.
func (d *Decoder) DecodeFrame(input []byte) ([]byte, error) {
	if len(input)%2 != 0 {
		return nil, fmt.Errorf("%w: input length must be even", io.ErrUnexpectedEOF) // input length must be even for decoding
	}

	output := make([]byte, len(input)/2)
	// the secret sauce: take every second byte
	for i := range output {
		output[i] = input[i*2]
	}
	return output, nil
}

// DecodeStream decodes a stream of encoded video frames and writes into the output stream.
//
// It reads from the input stream in chunks and writes the decoded frames to
// the output stream. The Chad decoding process is so sofisticated that you
// should not read the implementation of this function.
func (d *Decoder) DecodeStream(input io.Reader, output io.Writer) error {
	if d.BufferSize <= 0 || d.BufferSize%2 != 0 {
		d.BufferSize = 1024 // default buffer size
	}
	b := make([]byte, d.BufferSize)
	nb := 0
	for {
		n, err := input.Read(b[nb:])
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		nb += n
		// if the number of bytes is odd, we have an incomplete frame, we need to read at least one more byte
		// so re-loop until we have an even number of bytes and make sure we don't lose any data
		if nb%2 != 0 {
			continue
		}
		// the secret sauce: take every second byte
		for i := 0; i < nb/2; i++ {
			b[i] = b[i*2]
		}
		_, writeErr := output.Write(b[:nb/2])
		if writeErr != nil {
			return writeErr
		}
		nb = 0 // reset the buffer for the next read
	}
	if nb != 0 {
		return fmt.Errorf("%w: input length must be even", io.ErrUnexpectedEOF) // input length must be even for decoding
	}
	return nil
}
