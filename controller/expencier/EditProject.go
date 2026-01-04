package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditProjectBody struct {
	Name        string   `json:"name" validate:"required,min=3,max=50"`
	Description string   `json:"description" validate:"required,min=3,max=50"`
	SubLists    []string `json:"sub_lists"`
}

func (a *EditProjectBody) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	if len(validationErrors) == 0 {
		a.Name = strings.TrimSpace(a.Name)
		a.Description = strings.TrimSpace(a.Description)
	}
	return validationErrors
}

func EditProject(c *fiber.Ctx) error {
	auth := c.Locals("user").(*models.User)

	id := c.Params("id", "")
	idInObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	body := new(EditProjectBody)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	axFilter := bson.D{{Key: "_id", Value: idInObjectId}, {Key: "user_id", Value: auth.ID}}
	axResult, err := dbcontext.ExpencierProjectsModel.UpdateOne(c.Context(), axFilter, bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: body.Name}, {Key: "description", Value: body.Description}, {Key: "sub_lists", Value: body.SubLists}}}})

	bxFilter := bson.D{{Key: "project_id", Value: idInObjectId}, {Key: "user_id", Value: auth.ID}, {Key: "sub_list", Value: bson.D{{Key: "$nin", Value: body.SubLists}}}}
	bxResult, err := dbcontext.ExpencierTransactionsModel.UpdateMany(c.Context(), bxFilter, bson.D{{Key: "$set", Value: bson.D{{Key: "sub_list", Value: ""}}}})

	if axResult.ModifiedCount > 0 || bxResult.ModifiedCount > 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": false})
	}
}
