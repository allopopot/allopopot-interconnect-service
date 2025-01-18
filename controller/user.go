package controller

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/dto"
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

func (uc *UserController) SetPassword(c *fiber.Ctx) error {
	body := new(dto.SetPassword)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}
	user, _ := c.Locals("user").(*models.User)

	if !user.VerifyPassword(body.CurrentPassword) {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Failed to verify password."}})
	}
	user.SetPassword(body.SetPassword)
	dbResult, err := dbcontext.UserModel.UpdateByID(c.Context(), user.ID, bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: user.Password}}}})
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{err.Error()}})
	}

	if dbResult.ModifiedCount == 1 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": false})
	}
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	body := new(dto.UpdateUser)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}
	user, _ := c.Locals("user").(*models.User)

	if body.FirstName != "" {
		user.FirstName = body.FirstName
	}
	if body.Lastname != "" {
		user.LastName = body.Lastname
	}

	dbResult, err := dbcontext.UserModel.UpdateByID(c.Context(), user.ID, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{err.Error()}})
	}
	if dbResult.ModifiedCount == 1 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": false})
	}
}
