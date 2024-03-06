package routes

import (
	"rinha-backend-2024q1-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("alguem me da um emprego, é serio, EU PRECISO, aceito 600 reais")
	})
	app.Post("/clientes/:id/transacoes", controllers.CriarTransacao)

	extrato := app.Group("/clientes/:id/extrato")
	extrato.Get("/", controllers.ObterExtrato)
}
