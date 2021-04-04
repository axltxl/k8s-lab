package main

import (
	"log"

	"github.com/axltxl/k8s-lab/src/pkg/server"
)

const HTTP_PORT = "8000"

func main() {

	if err := server.Start(HTTP_PORT); err != nil {
		log.Fatal(err)
	}
}
