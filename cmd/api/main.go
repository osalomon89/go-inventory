package main

import (
	"log"

	"github.com/osalomon89/go-inventory/server"
)

func main() {
	const port string = ":8888"
	router := server.SetupRouter()

	log.Println("Server listening on port", port)

	err := server.Run(port, router)
	if err != nil {
		log.Fatalln(err)
	}
}
