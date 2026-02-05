package middleware

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/jsonwebtoken"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GateKeeper(c *fiber.Ctx) error {
	var token string = ""
	tokenFromCookie := c.Cookies("access_token")
	tokenFromAuthHeader := strings.Split(c.Get("authorization"), " ")

	if len(tokenFromCookie) > 0 {
		token = tokenFromCookie
	}

	if len(tokenFromAuthHeader) == 2 {
		token = tokenFromAuthHeader[1]
	}

	validatedToken, err := jsonwebtoken.ValidateToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		c.Cookie(&fiber.Cookie{
			Name:    "access_token",
			Value:   "",
			Expires: time.Now(),
		})
		c.Cookie(&fiber.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: time.Now(),
		})
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	if validatedToken.TokenType != "access" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Refresh Token Not Allowed   //////"},
		})
	}
	u := new(models.User)
	dbresult := dbcontext.UserModel.FindOne(context.TODO(), bson.D{{Key: "email", Value: validatedToken.PrincipalEmail}}).Decode(&u)
	if dbresult == mongo.ErrNoDocuments {
		c.Status(fiber.StatusUnauthorized)
		c.Cookie(&fiber.Cookie{
			Name:    "access_token",
			Value:   "",
			Expires: time.Now(),
		})
		c.Cookie(&fiber.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: time.Now(),
		})
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	c.Locals("user", u)
	return c.Next()
}
