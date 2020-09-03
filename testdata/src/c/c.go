package c

// S is a struct
type S struct{}

func (S) f() {}

func main() {
	var s S
	s.f() // want "NG"
}
