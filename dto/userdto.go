package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SetPassword struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	SetPassword     string `json:"set_password" validate:"required,min=8"`
}

func (a *SetPassword) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	return validationErrors
}

type UpdateUser struct {
	FirstName string `json:"first_name" validate:"omitempty,min=3"`
	Lastname  string `json:"last_name" validate:"omitempty,min=3"`
}

func (a *UpdateUser) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	return validationErrors
}
