package v1

import (
	"allopopot-interconnect-service/controller/expencier"
	"allopopot-interconnect-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func ExpencierRoutes(router fiber.Router) {

	router.Use(middleware.GateKeeper)

	router.Post("/project", expencier.CreateProject)
	router.Get("/project", expencier.GetProject)
	router.Delete("/project/:id", expencier.DeleteProject)
	router.Patch("/project/:id", expencier.EditProject)
}
