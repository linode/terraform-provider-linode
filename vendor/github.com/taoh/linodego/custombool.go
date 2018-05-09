package linodego

import (
	"fmt"
)

// CustomBool is a type to handle Linode's insistance of using ints as boolean values.
type CustomBool struct {
	Bool bool
}

func (cb *CustomBool) UnmarshalJSON(b []byte) error {
	if len(b) != 1 {
		return fmt.Errorf("Unable to marshal value with length %d into a CustomBool.", len(b))
	}
	if int(b[0]) == 0 {
		cb.Bool = false
	} else if int(b[0]) == 1 {
		cb.Bool = true
	}
	return nil
}

func (cb *CustomBool) MarshalJSON() ([]byte, error) {
	if cb == nil {
		return []byte{}, fmt.Errorf("Unable to marshal nil value for CustomBool")
	}
	if cb.Bool {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}
