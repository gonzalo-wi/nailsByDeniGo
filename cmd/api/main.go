package main

import (
	"log"

	bootstrap "apiGoShei/internal/boostrap"
)

func main() {
	app, port := bootstrap.BuildApp()
	log.Printf("Servidor iniciando en :%s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
