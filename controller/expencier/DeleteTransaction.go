package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTransaction(c *fiber.Ctx) error {

	auth := c.Locals("user").(*models.User)

	pid := c.Params("pid")
	tid := c.Params("tid")

	projectId, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Cannot parse Product ID"}})
	}
	transactionId, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Cannot parse Transaction ID"}})
	}

	axResult, err := dbcontext.ExpencierTransactionsModel.DeleteOne(c.Context(), bson.D{{Key: "_id", Value: transactionId}, {Key: "project_id", Value: projectId}, {Key: "user_id", Value: auth.ID}})
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not delete transaction"}})
	}

	if axResult.DeletedCount > 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": false})
	}

}
