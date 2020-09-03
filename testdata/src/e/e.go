package e

func main() {
	err := func() {}
	err() // want "NG"

	func() {}() // want "NG"
}
