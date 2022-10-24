package main

//TODO: Factorización y formateo de errores
//TODO: Consultas con orm
//TODO: Realizaar cambios para un nuevo PR
import (
	"log"

	"github.com/osalomon89/go-inventory/internal/server"
)

func main() {

	const port string = ":8888"
	httpServer := server.NewHTTPRouter(port)

	router := httpServer.SetupRouter()

	log.Println("Server listening on port", port)

	err := httpServer.Run(router)
	if err != nil {
		log.Fatalln(err)
	}
}
