package middleware

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/jsonwebtoken"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GateKeeper(c *fiber.Ctx) error {
	token := strings.Split(c.Get("authorization"), " ")[1]
	validatedToken, err := jsonwebtoken.ValidateToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	u := new(models.User)
	dbresult := dbcontext.DB.Where("email = ?", validatedToken.PrincipalEmail).First(&u)
	if dbresult.RowsAffected == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": []string{"//////   GateKeeper Says: Not Authorized   //////"},
		})
	}
	c.Locals("user", u)
	return c.Next()
}
