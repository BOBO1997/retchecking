package b

func f() {}

func main() {
	f() // want "NG"

	if true {
		f()
		return
	}

	f() // want "NG"
}
