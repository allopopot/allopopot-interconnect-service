package controller

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/dto"
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/jsonwebtoken"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct{}

func (ac *AuthController) CreateAccount(c *fiber.Ctx) error {
	body := new(dto.CreateAccount)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	u := new(models.User)
	u.Email = body.Email

	result := dbcontext.DB.First(&u, "email = ?", u.Email)
	if result.RowsAffected > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Account Already Exists"}})
	}
	u.FirstName = body.FirstName
	u.LastName = body.LastName
	u.SetPassword(body.Password)
	result = dbcontext.DB.Create(u)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{result.Error.Error()}})
	}
	return c.JSON(fiber.Map{"data": u})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	body := new(dto.Login)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	u := new(models.User)
	result := dbcontext.DB.First(u, "email = ?", body.Email)
	if result.RowsAffected == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"error": []string{"Please check your credentials."}})
	}

	ok := u.VerifyPassword(body.Password)
	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"error": []string{"Please check your credentials."}})
	}

	claims := new(jsonwebtoken.AIMClaims)
	claims.ID = string(rune(u.ID))
	claims.PrincipalID = int(u.ID)
	claims.PrincipalName = fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	claims.PrincipalEmail = u.Email
	signedString, err := jsonwebtoken.GenerateToken(*claims)

	if err != nil {
		log.Panicln("Cannot Generate JWT", err)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{"access_token": signedString}})
}

func (ac *AuthController) VerifyToken(c *fiber.Ctx) error {
	token := c.Get("authorization")
	if len(token) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Authorization Header Invalid"}})
	}
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Authorization Header Invalid"}})
	}
	claims, err := jsonwebtoken.ValidateToken(tokenParts[1])
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid Token"}})
	}
	return c.JSON(fiber.Map{"data": claims})
}
