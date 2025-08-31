package v1

import (
	"allopopot-interconnect-service/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	ac := new(controller.AuthController)
	router.Post("/createAccount", ac.CreateAccount)
	router.Post("/login", ac.Login)
	router.Get("/verifyToken", ac.VerifyToken)
	router.Post("/refreshToken", ac.RefreshToken)
	router.Post("/logout", ac.Logout)
}
