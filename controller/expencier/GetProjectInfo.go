package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProjectInfo(c *fiber.Ctx) error {
	id := c.Params("pid", "")

	var filter bson.D
	filter = append(filter, bson.E{Key: "user_id", Value: c.Locals("user").(*models.User).ID})

	if id != "" {
		idInObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Cannot parse ID"}})
		}
		filter = append(filter, bson.E{Key: "project_id", Value: idInObjectId})
	}

	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: filter},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "type",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									bson.D{
										{Key: "$gt",
											Value: bson.A{
												"$amount",
												0,
											},
										},
									},
									"income",
									"expense",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$type"},
					{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
					{Key: "count", Value: bson.D{{Key: "$count", Value: bson.D{}}}},
				},
			},
		},
		bson.D{{
			Key: "$set", Value: bson.D{{Key: "type", Value: "$_id"}},
		}},
		bson.D{{
			Key: "$unset", Value: "_id",
		}},
	}

	cursor, err := dbcontext.ExpencierTransactionsModel.Aggregate(c.Context(), pipeline)
	if err != nil {
		println(err.Error())
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not get transactions in the project"}})
	}
	defer cursor.Close(c.Context())

	var results []bson.M
	cursor.All(c.Context(), &results)

	return c.JSON(fiber.Map{"data": results})
}
