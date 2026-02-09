package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditTransactionBody struct {
	SubList     string  `json:"sub_list"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description" validate:"required,min=3,max=100"`
	EntryDate   string  `json:"entry_date" validate:""`
}

func (a *EditTransactionBody) Validate() []string {
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

func EditTransaction(c *fiber.Ctx) error {

	auth := c.Locals("user").(*models.User)

	body := new(EditTransactionBody)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	var entryDate time.Time
	if body.EntryDate != "" {
		ed, err := time.Parse(time.RFC3339, body.EntryDate)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Could not parse date"}})
		} else {
			entryDate = ed
		}
	}

	id := c.Params("pid")
	projectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	tid := c.Params("tid")
	transactionID, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	filter := bson.D{{Key: "_id", Value: transactionID}, {Key: "project_id", Value: projectID}, {Key: "user_id", Value: auth.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "amount", Value: body.Amount}, {Key: "description", Value: body.Description}, {Key: "sub_list", Value: body.SubList}, {Key: "created_time", Value: entryDate}}}}

	axResult, err := dbcontext.ExpencierTransactionsModel.UpdateOne(c.Context(), filter, update)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not edit transaction"}})
	}

	if axResult.ModifiedCount > 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"data": true})
	} else {
		c.Status(fiber.StatusNotModified)
		return c.JSON(fiber.Map{"data": false})
	}

}
