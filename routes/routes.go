package routes

import (
	"rinha-backend-2024q1-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/clientes/:id/transacoes", controllers.ProcessarTransacao)

	extrato := app.Group("/clientes/:id/extrato")
	extrato.Get("/", controllers.ObterExtrato)
}
