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

type CreateProjectBody struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3,max=50"`
}

func (a *CreateProjectBody) Validate() []string {
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

func CreateProject(c *fiber.Ctx) error {

	auth := c.Locals("user").(*models.User)

	body := new(CreateProjectBody)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	project := new(models.ExpencierProjects)
	project.ID = primitive.NewObjectID()
	project.Name = body.Name
	project.Description = body.Description
	project.UserId = auth.ID
	project.CreatedTime = time.Now()

	_, err := dbcontext.ExpencierProjectsModel.InsertOne(c.Context(), project)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{err.Error()}})
	}

	return c.JSON(fiber.Map{"data": project})
}
