package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/henriiquematheus/rinha-backend-2024q1-go/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/clientes/:id/transacoes", controllers.ProcessarTransacao)
	// Rotas relacionadas a extrato
	extrato := app.Group("/clientes/:id/extrato")
	extrato.Get("/", controllers.ObterExtrato)
}