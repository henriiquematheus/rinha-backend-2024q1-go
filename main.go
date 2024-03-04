// main.go
package main

import (
	"log"
	connection "rinha-backend-2024q1-go/db"
	"rinha-backend-2024q1-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	connectionString := "postgres://postgres:postgres@db:5432/rinha"

	err := connection.Init(connectionString)
	if err != nil {
		log.Fatal("Erro ao inicializar a conexão com o banco de dados:", err)
	}

	defer connection.GetPool().Close()

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Println("Servidor iniciado na porta 8080")

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
