# Chad encoding algorithm (Go)

The Chad encoding algorithm can be used as a library.

```
go get github.com/luisferreira32/chadenc
```

There are two modes:

- Atomic encoding/decoding of single frames
- Streaming mode

## Atomic encoding/decoding of single frames

Assuming rendering or source of raw video frames.

```go
// encode on the server the frame
encodedFrame, err := chadenc.EncodeFrame(frame)
if err != nil {
    // error handling
}
// send encoded frame down the video stream

// on the client side decode the frame
decodedFrame, err := chadenc.DecodeFrame()
if err != nil {
    // error handling
}
// display the decodedFrame
```

## Streaming mode

```go
// any io.Reader / io.Writer streams
encoder := &chadenc.Encoder{}
err := encoder.EncodeStream(inputStream, outputStream)
if err != nil {
    // error handling
}

// on the client side of the stream
decoder := &chadenc.Decoder{}
err := decoder.DecodeStream(encodedStream, displayStream)
if err != nil {
    // error handling
}
// displayStream must be broken into displayable frames
```
