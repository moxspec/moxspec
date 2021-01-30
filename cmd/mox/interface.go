package main

// Decoder is implemented by any value that provides information for mox
type Decoder interface {
	Decode() error
}
