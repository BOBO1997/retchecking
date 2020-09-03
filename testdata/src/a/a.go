package a

import (
	"fmt"
)

// S is a fake
type S struct{}

func (S) Error() error {
	return fmt.Errorf("this is error")
}

func fError() error {
	return fmt.Errorf("this is error")
}

func main() {
	if true {
		fError()
		return // OK
	}

	var s S
	s.Error() // want "NG"

	a := 1
	if a == 1 {
		s.Error()
		return // OK
	}

	fError() // want "NG"

	f := func(a int) error { return fmt.Errorf("%d, this is also error", a) }
	f(a) // want "NG"
}
