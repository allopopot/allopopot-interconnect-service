package middleware

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/jsonwebtoken"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GateKeeper(c *fiber.Ctx) error {
	tokenParts := strings.Split(c.Get("authorization"), " ")
	if len(tokenParts) != 2 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	token := tokenParts[1]
	validatedToken, err := jsonwebtoken.ValidateToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
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
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	c.Locals("user", u)
	return c.Next()
}
