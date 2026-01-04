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

type DeleteTransactionBody struct {
	SubList     string  `json:"sub_list"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description" validate:"required,min=3,max=100"`
}

func (a *DeleteTransactionBody) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	if len(validationErrors) == 0 {
		a.Description = strings.TrimSpace(a.Description)
		a.SubList = strings.TrimSpace(a.SubList)
	}
	return validationErrors
}

func DeleteTransaction(c *fiber.Ctx) error {

	auth := c.Locals("user").(*models.User)

	body := new(CreateTransactionBody)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

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
