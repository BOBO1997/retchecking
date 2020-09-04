package e

import "net/http"

// S is a struct
type S struct {
	f func(http.ResponseWriter, string, int)
}

func main() {
	err1, err2 := func() {}, func() {}
	err1() // want "NG"
	err2() // want "NG"

	func() {}() // want "OK"

	var err3 func()
	err3 = func() {}

	err3() // want "NG"

	// how about []func() list?
	errs := make([]func(), 10)
	errs[5] = func() {}
	errs[5]() // want "NG"

	var s S
	var w http.ResponseWriter
	s.f = http.Error

	s.f(w, "", 0) // want "NG"
}
