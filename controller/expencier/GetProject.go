package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProject(c *fiber.Ctx) error {
	id := c.Query("pid", "")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)

	if limit > 20 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Limit cannot be greater than 20"}})
	}

	var filter bson.D
	filter = append(filter, bson.E{Key: "user_id", Value: c.Locals("user").(*models.User).ID})

	if id != "" {
		idInObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Cannot parse ID"}})
		}
		filter = append(filter, bson.E{Key: "_id", Value: idInObjectId})
	}

	cursor, err := dbcontext.ExpencierProjectsModel.Find(c.Context(), filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not transactions in the project"}})
	}
	defer cursor.Close(c.Context())

	var results []models.ExpencierProjects
	// var results []bson.M   // WARNING: ALTERNATE TO ABOVE STATEMENT
	cursor.All(c.Context(), &results)

	return c.JSON(fiber.Map{"data": results})
}
