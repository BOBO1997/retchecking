package e

import "net/http"

// S is a struct
type S struct {
	f func(http.ResponseWriter, string, int)
}

func main() {
	var s S
	var w http.ResponseWriter

	err1, err2 := http.Error, http.Error
	err1(w, "", 0) // want "NG"
	err2(w, "", 0) // want "NG"

	func() {}() // want "OK"

	var err3 func()
	err3 = func() {}

	err3() // want "NG"

	// how about []func() list?
	errs := make([]func(), 10)
	// ! this cannot be resolved
	errs[5] = func() {}
	errs[5]() // want "NG"

	s.f = http.Error

	s.f(w, "", 0) // want "NG"
}
