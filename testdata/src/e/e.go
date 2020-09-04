package e

func main() {
	err1, err2 := func() {}, func() {}
	err1() // want "NG"
	err2() // want "NG"

	func() {}() // want "OK"
}
