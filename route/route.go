package route

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API Running 🚀",
		})
	})

	// TODO: Auth routes
	// TODO: User routes
	// TODO: Achievement routes
}
