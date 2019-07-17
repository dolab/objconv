package resp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dolab/objconv"
)

func TestParser(t *testing.T) {
	for _, test := range respDecodeTests {
		t.Run(testName(test.s), func(t *testing.T) {
			r := bytes.NewReader([]byte(test.s))
			p := NewParser(r)

			typ, err := p.ParseType()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(test.t, typ) {
				t.Errorf("expected %#v with %s, but got %v", test, test.t.String(), typ.String())
			}
		})
	}
}

func TestParserWithPipeline(t *testing.T) {
	cmds := "*5\r\n$3\r\nset\r\n$3\r\nk-0\r\n$112\r\nresp: ParseMapNext should never be called because RESP has no map type, this is likely a bug in the decoder code\r\n$2\r\nex\r\n$1\r\n1\r\n*2\r\n$3\r\nget\r\n$3\r\nk-0\r\n*5\r\n$3\r\nset\r\n$3\r\nk-1\r\n$1\r\n1\r\n$2\r\nex\r\n$1\r\n1\r\n*2\r\n$3\r\nget\r\n$3\r\nk-1\r\n"
	pipe := bytes.NewReader([]byte(cmds))

	p := NewParser(pipe)

	// ParseType
	typ, err := p.ParseType()
	if err != nil {
		t.Error(err)
	}

	if typ != objconv.Array {
		t.Errorf("expect message type of %s, but got %s", objconv.Array, typ)
	}

	// read array
	n, err := p.ParseArrayBegin()
	if err != nil {
		t.Error(err)
	}

	arrs := []string{
		"set",
		"k-0",
		"resp: ParseMapNext should never be called because RESP has no map type, this is likely a bug in the decoder code",
		"ex",
		"1",
	}
	for i := 0; i < n; i++ {
		b, err := p.ParseBytes()
		if err != nil {
			t.Error(err)
		}

		if arrs[i] != string(b) {
			t.Errorf("expect bytes of %q, but got %q", arrs[i], string(b))
		}
	}

	// ParseType again
	typ, err = p.ParseType()
	if err != nil {
		t.Error(err)
	}

	if typ != objconv.Array {
		t.Errorf("expect message type of %s, but got %s", objconv.Array, typ)
	}

	// read array again
	n, err = p.ParseArrayBegin()
	if err != nil {
		t.Error(err)
	}

	arrs = []string{
		"get",
		"k-0",
	}
	for i := 0; i < n; i++ {
		b, err := p.ParseBytes()
		if err != nil {
			t.Error(err)
		}

		if arrs[i] != string(b) {
			t.Errorf("expect bytes of %q, but got %q", arrs[i], string(b))
		}
	}
}

func TestParserWithPipelineOfStreamDecoder(t *testing.T) {
	cmds := "*5\r\n$3\r\nset\r\n$3\r\nk-0\r\n$112\r\nresp: ParseMapNext should never be called because RESP has no map type, this is likely a bug in the decoder code\r\n$2\r\nex\r\n$1\r\n1\r\n*2\r\n$3\r\nget\r\n$3\r\nk-0\r\n*5\r\n$3\r\nset\r\n$3\r\nk-1\r\n$1\r\n1\r\n$2\r\nex\r\n$1\r\n1\r\n*2\r\n$3\r\nget\r\n$3\r\nk-1\r\n"
	pipe := bytes.NewReader([]byte(cmds))

	p := NewParser(pipe)

	dec := objconv.NewStreamDecoder(p)

	// first cmd
	arrs := []string{
		"set",
		"k-0",
		"resp: ParseMapNext should never be called because RESP has no map type, this is likely a bug in the decoder code",
		"ex",
		"1",
	}

	if n := dec.Len(); n != len(arrs) {
		t.Error("invalid length returned by the stream decoder:", n)
	}

	var (
		v []byte
		i int
	)
	for dec.Decode(&v) == nil {
		if arrs[i] != string(v) {
			t.Error(string(v), "!=", arrs[i])
		}
		i++
	}

	// reset for piped command
	if err := dec.Next(); err != nil {
		t.Error(err)
	}

	// second cmd
	arrs = []string{
		"get",
		"k-0",
	}

	if n := dec.Len(); n != len(arrs) {
		t.Error("invalid length returned by the stream decoder:", n)
	}

	v = []byte{}
	i = 0
	for dec.Decode(&v) == nil {
		if arrs[i] != string(v) {
			t.Error(string(v), "!=", arrs[i])
		}
		i++
	}

	// reset for piped command
	if err := dec.Next(); err != nil {
		t.Error(err)
	}

	// third cmd
	arrs = []string{
		"set",
		"k-1",
		"1",
		"ex",
		"1",
	}

	if n := dec.Len(); n != len(arrs) {
		t.Error("invalid length returned by the stream decoder:", n)
	}

	v = []byte{}
	i = 0
	for dec.Decode(&v) == nil {
		if arrs[i] != string(v) {
			t.Error(string(v), "!=", arrs[i])
		}
		i++
	}

	// reset for piped command
	if err := dec.Next(); err != nil {
		t.Error(err)
	}

	// fourth cmd
	arrs = []string{
		"get",
		"k-1",
	}

	if n := dec.Len(); n != len(arrs) {
		t.Error("invalid length returned by the stream decoder:", n)
	}

	v = []byte{}
	i = 0
	for dec.Decode(&v) == nil {
		if arrs[i] != string(v) {
			t.Error(string(v), "!=", arrs[i])
		}
		i++
	}

	// end of parser
	if err := dec.Decode(&v); err != objconv.End {
		t.Errorf("%v != %v", err, objconv.End)
	}
}
