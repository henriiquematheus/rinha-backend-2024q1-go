// main.go
package main

import (
	"log"

	"rinha-backend-2024q1-go/db"
	"github.com/gofiber/fiber/v2"
	"rinha-backend-2024q1-go/routes"
)

func main() {
	// Substitua a string de conexão conforme necessário
	connectionString := "postgres://postgres:postgres@db:5432/rinha"

	err := connection.Init(connectionString)
	if err != nil {
		log.Fatal("Erro ao inicializar a conexão com o banco de dados:", err)
	}

	defer connection.GetPool().Close()

	app := fiber.New()
	routes.SetupRoutes(app)

	// Adicionando mensagem de log
	log.Println("Servidor iniciado na porta 8080")

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
