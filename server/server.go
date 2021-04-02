package server

import (
	"fmt"
	"net/http"
)

func Start(port string) error {

	fmt.Printf("Starting server at port %s", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
