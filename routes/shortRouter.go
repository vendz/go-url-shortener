package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/vendz/bitsnip/controllers"
)

func ShortRoutes(shortRoutes *fiber.App, h *controllers.Database) {
	shortRoutes.Get("/:alias", h.GetUrl)
	shortRoutes.Post("/", h.ShortenUrl)
	shortRoutes.Delete("/:alias", h.DeleteUrl)
}
