package resp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dolab/objconv"
)

var respDecodeTests = []struct {
	v interface{}
	s string
	t objconv.Type
}{
	{nil, "$-1\r\n", objconv.Nil},
	{nil, "*-1\r\n", objconv.Nil},

	{0, ":0\r\n", objconv.Int},
	{-1, ":-1\r\n", objconv.Int},
	{42, ":42\r\n", objconv.Int},

	{"", "+\r\n", objconv.String},
	{"Hello World!", "+Hello World!\r\n", objconv.String},
	{"Hello\nWorld!", "+Hello\nWorld!\r\n", objconv.String},
	{"Hello\r\nWorld!", "$13\r\nHello\r\nWorld!\r\n", objconv.Bytes},

	{[]byte{}, "$0\r\n\r\n", objconv.Bytes},
	{[]byte("Hello World!"), "$12\r\nHello World!\r\n", objconv.Bytes},

	{NewError(""), "-\r\n", objconv.Error},
	{NewError("oops"), "-oops\r\n", objconv.Error},
	{NewError("ERR A"), "-ERR A\r\n", objconv.Error},

	{[]int{}, "*0\r\n", objconv.Array},
	{[]int{1, 2, 3}, "*3\r\n:1\r\n:2\r\n:3\r\n", objconv.Array},
}

func TestUnmarshal(t *testing.T) {
	for _, test := range respDecodeTests {
		t.Run(testName(test.s), func(t *testing.T) {
			var typ reflect.Type

			if test.v == nil {
				typ = reflect.TypeOf((*interface{})(nil)).Elem()
			} else {
				typ = reflect.TypeOf(test.v)
			}

			val := reflect.New(typ)
			err := Unmarshal([]byte(test.s), val.Interface())

			if err != nil {
				t.Error(err)
			}

			v1 := test.v
			v2 := val.Elem().Interface()

			if !reflect.DeepEqual(v1, v2) {
				t.Error(v2)
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, test := range respDecodeTests {
		var t reflect.Type

		if test.v == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test.v)
		}

		v := reflect.New(t).Interface()
		s := []byte(test.s)

		b.Run(testName(test.s), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := Unmarshal(s, v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}

func BenchmarkDecoder(b *testing.B) {
	r := strings.NewReader("")
	p := NewParser(nil)
	d := objconv.NewDecoder(p)

	for _, test := range respDecodeTests {
		var t reflect.Type

		if test.v == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test.v)
		}

		v := reflect.New(t).Interface()

		b.Run(testName(test.s), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				r.Reset(test.s)
				p.Reset(r)

				if err := d.Decode(v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}
