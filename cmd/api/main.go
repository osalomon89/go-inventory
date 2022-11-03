// Inventory APP.
//
// API to create, update and delete books.
//
//	 Schemes: http, https
//	 Host: localhost:5000
//		BasePath: /
//		Version: 1.0
//		License: MIT http://opensource.org/licenses/MIT
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"log"

	"github.com/osalomon89/go-inventory/internal/server"
)

func main() {
	const port string = ":8888"
	httpServer := server.NewHTTPRouter(port)

	router := httpServer.SetupRouter()

	log.Println("Server listening on port", port)

	if err := httpServer.Run(router); err != nil {
		log.Fatalln(err)
	}
}
