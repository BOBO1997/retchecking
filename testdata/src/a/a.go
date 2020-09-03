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

func g() {}

func h() {}

func main() {
	if true {
		fError() // want "OK"
		return
	}

	var s S
	s.Error() // want "NG"

	a := 1
	if a == 1 {
		s.Error() // want "OK"
		return
	}

	fError() // want "NG"

	f := func(a int) error { return fmt.Errorf("%d, this is also error", a) }
	f(a) // want "NG"

	serror := s.Error
	serror() // want "NG"

	{
		g() // want "OK"

		if true {
			s.Error() // want "OK"
			return
		}
	}
	h() // want "OK"

}

/*
OK
NG
OK
NG
NG
NG
OK
OK
OK
*/
