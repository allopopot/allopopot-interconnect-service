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
	app.Use(cors.New())

	app.Route("/v1", routes.V1)
	app.Listen(config.SERVER_URI)
}
