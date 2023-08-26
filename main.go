package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vendz/bitsnip/database"
	"github.com/vendz/bitsnip/helper"
	"github.com/vendz/bitsnip/routes"
)

func main() {
	helper.LoadEnv()
	handler := database.NewDatabase()
	defer handler.RedisClient.Close()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 30 * time.Second,
	}))
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is Up and Running...")
	})
	routes.ShortRoutes(app, &handler)

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
