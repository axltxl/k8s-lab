package main

import (
	"log"

	"github.com/axltxl/k8s-lab/src/pkg/server"
)

func main() {

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
