package main

import (
	"github.com/cripplet/clicker/clicker-rest/lib"
	"net/http"
)

func main() {
	http.HandleFunc("/", cc_rest_lib.GameRouter)
	http.ListenAndServe(":8080", nil)
}
