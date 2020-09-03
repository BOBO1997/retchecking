package d

import "net/http"

func main() {
	var w http.ResponseWriter
	http.Error(w, "aaaaa", 0) // want "OK"
	return
}
