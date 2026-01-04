package routes

import (
	v1 "allopopot-interconnect-service/routes/v1"

	"github.com/gofiber/fiber/v2"
)

func V1(router fiber.Router) {
	router.Route("/auth", v1.AuthRoutes)
	router.Route("/user", v1.UserRoutes)
	router.Route("/expencier", v1.ExpencierRoutes)
}
