package main

import (
	"log"

	"github.com/osalomon89/go-inventory/internal/server"
)

func main() {
	const port string = ":8080"
	httpServer := server.NewHTTPRouter(port)

	router := httpServer.SetupRouter()

	log.Println("Server listening on port", port)

	if err := httpServer.Run(router); err != nil {
		log.Fatalln(err)
	}
}
