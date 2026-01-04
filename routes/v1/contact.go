package v1

import (
	"allopopot-interconnect-service/controller/contact"
	"allopopot-interconnect-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ContactRoutes(router fiber.Router) {

	router.Use(middleware.GateKeeper)

	router.Post("/contact", contact.SearchContact)
}
