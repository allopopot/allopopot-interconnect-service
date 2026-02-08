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
	router.Get("/project/:pid", expencier.GetProjectInfo)
	router.Delete("/project/:pid", expencier.DeleteProject)
	router.Patch("/project/:pid", expencier.EditProject)
	router.Post("/project/:pid/transaction", expencier.CreateTransaction)
	router.Get("/project/:pid/transaction", expencier.GetTransactions)
	router.Delete("/project/:pid/transaction/:tid", expencier.DeleteTransaction)
	router.Patch("/project/:pid/transaction/:tid", expencier.EditTransaction)
}
