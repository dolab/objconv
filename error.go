package objconv

import (
	"errors"
	"fmt"
)

var (
	// End is expected to be returned to indicate that a function has completed
	// its work, this is usually employed in generic algorithms.
	End = errors.New("end")

	// Shadow is expected to be returned to indicate that a parser has completed
	// its work, but there is extra bytes left in the buffer.
	Shadow = errors.New("shadow")
)

func typeConversionError(from Type, to Type) error {
	return fmt.Errorf("objconv: cannot convert from %s to %s", from, to)
}
