package main

import (
	"log"

	"github.com/osalomon89/go-inventory/internal/server"
)

func main() {

	const port string = ":8888" // se crea una constante para fijar la la direƒccion de la ruta

	//reemplazamos la variable que creamos de NewRouter por:

	httpServer := server.NewHTTPRouter(port)
	router := httpServer.SetupRouter() // exportar las diferentes rutas que estan en el modulo de server a traves de la asignacion de variable

	log.Println("Server listening in port:", port)

	if err := httpServer.Run(router); err != nil {
		log.Fatal(err) // se controla el error
	}

}
