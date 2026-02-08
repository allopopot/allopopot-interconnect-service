package expencier

import (
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/models"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTransactionBody struct {
	SubList     string  `json:"sub_list"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description" validate:"required,min=3,max=100"`
	EntryDate   string  `json:"entry_date" validate:""`
}

func (a *CreateTransactionBody) Validate() []string {
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

func CreateTransaction(c *fiber.Ctx) error {

	auth := c.Locals("user").(*models.User)

	body := new(CreateTransactionBody)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	id := c.Params("pid")
	projectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not parse project ID"}})
	}

	project := new(models.ExpencierProjects)
	axresult := dbcontext.ExpencierProjectsModel.FindOne(c.Context(), primitive.D{{Key: "_id", Value: projectID}, {Key: "user_id", Value: auth.ID}}).Decode(&project)

	if axresult != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not find project"}})
	}

	newTransaction := new(models.ExpencierTransactions)
	newTransaction.ID = primitive.NewObjectID()
	newTransaction.ProjectId = projectID
	newTransaction.UserID = auth.ID
	newTransaction.Amount = body.Amount
	newTransaction.Description = body.Description
	newTransaction.SubList = body.SubList
	newTransaction.CreatedTime = time.Now()

	if body.EntryDate != "" {
		newTransaction.CreatedTime, err = time.Parse(time.RFC3339, body.EntryDate)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Could not parse date"}})
		}
	}

	_, err = dbcontext.ExpencierTransactionsModel.InsertOne(c.Context(), newTransaction)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Could not create transaction"}})
	}

	return c.JSON(fiber.Map{"data": newTransaction})
}
