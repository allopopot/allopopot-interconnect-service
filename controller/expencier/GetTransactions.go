package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTransactions(c *fiber.Ctx) error {
	id := c.Query("pid", "")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)
	startdate := c.Query("startdate", "")
	enddate := c.Query("enddate", "")

	startDate, err := time.Parse(time.DateOnly, startdate)
	endDate, err := time.Parse(time.DateOnly, enddate)

	auth := c.Locals("user").(*models.User)

	pid := c.Params("pid")
	projectID, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	if limit > 100 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Limit cannot be greater than 20"}})
	}

	var filter bson.D
	filter = append(filter, bson.E{Key: "user_id", Value: auth.ID})
	filter = append(filter, bson.E{Key: "project_id", Value: projectID})
	if startdate != "" {
		filter = append(filter, bson.E{Key: "created_time", Value: bson.D{{Key: "$gte", Value: startDate}}})
	}
	if enddate != "" {
		filter = append(filter, bson.E{Key: "created_time", Value: bson.D{{Key: "$lte", Value: endDate}}})
	}

	if id != "" {
		idInObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Cannot parse ID"}})
		}
		filter = append(filter, bson.E{Key: "_id", Value: idInObjectId})
	}

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: filter}},
		// bson.D{
		// 	{Key: "$lookup",
		// 		Value: bson.D{
		// 			{Key: "from", Value: "expencier_projects"},
		// 			{Key: "localField", Value: "project_id"},
		// 			{Key: "foreignField", Value: "_id"},
		// 			{Key: "as", Value: "project"},
		// 			{Key: "pipeline",
		// 				Value: bson.A{
		// 					bson.D{
		// 						{Key: "$unset",
		// 							Value: bson.A{
		// 								"sub_lists",
		// 								"user_id",
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// bson.D{{Key: "$set", Value: bson.D{{Key: "project", Value: bson.D{{Key: "$first", Value: "$project"}}}}}},
		bson.D{
			{Key: "$facet",
				Value: bson.D{
					{Key: "data",
						Value: bson.A{
							bson.D{{Key: "$skip", Value: skip}},
							bson.D{{Key: "$limit", Value: limit}},
						},
					},
					{Key: "count",
						Value: bson.A{
							bson.D{{Key: "$count", Value: "count"}},
						},
					},
				},
			},
		},
		bson.D{{Key: "$set", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$first", Value: "$count.count"}}}}}},
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

	return c.JSON(results[0])

}
