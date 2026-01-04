package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("pid", "")

	auth := c.Locals("user").(*models.User)

	idInObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	ax, err := dbcontext.ExpencierProjectsModel.DeleteOne(c.Context(), bson.D{{Key: "_id", Value: idInObjectId}, {Key: "user_id", Value: auth.ID}})
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not delete project"}})
	}

	bx, err := dbcontext.ExpencierTransactionsModel.DeleteMany(c.Context(), bson.D{{Key: "project_id", Value: idInObjectId}, {Key: "user_id", Value: auth.ID}})
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not transactions in the project"}})
	}

	if ax.DeletedCount > 0 || bx.DeletedCount > 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": "Could not delete project"})
	}
}
