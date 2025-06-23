package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CreateAccount struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (a *CreateAccount) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	if len(validationErrors) == 0 {
		a.FirstName = strings.TrimSpace(a.FirstName)
		a.LastName = strings.TrimSpace(a.LastName)
		a.Email = strings.TrimSpace(a.Email)
		a.Email = strings.ToLower(a.Email)
	}
	return validationErrors
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (a *Login) Validate() []string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(a)
	var validationErrors []string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}
	if len(validationErrors) == 0 {
		a.Email = strings.TrimSpace(a.Email)
		a.Email = strings.ToLower(a.Email)
	}
	return validationErrors
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=8"`
}

func (a *RefreshToken) Validate() []string {
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
