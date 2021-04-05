package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/axltxl/k8s-lab/src/pkg/config"
)

// FIXME: doc me
// var config.HttpPort string = config.Get().HttpPort
// var config.HttpPort string = config.HttpPort

/* todoList resource */
func todoListHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello todoListHandler!")
}

/* todoListItem resource */
func todoListItemHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello todoListItemHandler!")
}

// FIXME: doc me
func setupHandlers() {
	http.HandleFunc("/todolist", todoListHandler)
	http.HandleFunc("/todolist/item", todoListItemHandler)
}

func Start() error {

	//
	setupHandlers()

	// FIXME: doc me
	log.Printf("Starting server at port %s", config.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.HttpPort), nil)
}
