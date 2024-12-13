package controller

import (
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (uc *UserController) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Cannot Parse User Model"}})
	}
	return c.JSON(fiber.Map{"data": user})
}
