package resp

import (
	"bytes"
	"reflect"
	"testing"
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
