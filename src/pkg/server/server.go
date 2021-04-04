package server

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func Start(port string) error {

	fmt.Printf("Starting server at port %s", port)

	//
	http.HandleFunc("/hello", helloHandler)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
