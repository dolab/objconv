package yaml

import (
	"testing"

	"github.com/dolab/objconv/objtests"
)

func TestCodec(t *testing.T) {
	objtests.TestCodec(t, Codec)
}

func BenchmarkCodec(b *testing.B) {
	objtests.BenchmarkCodec(b, Codec)
}
