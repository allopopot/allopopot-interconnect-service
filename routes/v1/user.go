package v1

import (
	"allopopot-interconnect-service/controller"
	"allopopot-interconnect-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	uc := new(controller.UserController)

	router.Use(middleware.GateKeeper)

	router.Get("/me", uc.Me)
	router.Post("/setPassword", uc.SetPassword)
	router.Post("/updateUser", uc.UpdateUser)
}
