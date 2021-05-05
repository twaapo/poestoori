package main

import (
	"fmt"
	"net/http"
	"poegen"
)

func serverecipe(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, poegen.Generate())
}

func main() {
	http.HandleFunc("/", serverecipe)
	http.ListenAndServe(":8090", nil)
}
