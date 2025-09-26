package main

import (
	"allopopot-interconnect-service/config"
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/routes"
	"allopopot-interconnect-service/service/emailqueue"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	dbcontext.InitDb()
	emailqueue.QueueConnect()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
	}))
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     strings.Join(config.CORS_ORIGINS, ","),
	// 	AllowCredentials: true,
	// }))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})
	app.Route("/v1", routes.V1)
	app.Listen(config.SERVER_URI)
}
