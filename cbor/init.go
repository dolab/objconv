package cbor

import (
	"io"

	"github.com/dolab/objconv"
)

// Codec for the CBOR format.
var Codec = objconv.Codec{
	NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
	NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
}

func init() {
	for _, name := range [...]string{
		"application/cbor",
		"cbor",
	} {
		objconv.Register(name, Codec)
	}
}
