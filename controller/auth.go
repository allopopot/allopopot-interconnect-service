package controller

import (
	"allopopot-interconnect-service/config"
	"allopopot-interconnect-service/dbcontext"
	"allopopot-interconnect-service/dto"
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/emailqueue"
	"allopopot-interconnect-service/service/emailtemplates"
	"allopopot-interconnect-service/service/jsonwebtoken"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct{}

func (ac *AuthController) CreateAccount(c *fiber.Ctx) error {
	body := new(dto.CreateAccount)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	u := new(models.User)
	u.Email = body.Email

	resultA := dbcontext.UserModel.FindOne(c.Context(), bson.D{{Key: "email", Value: u.Email}}).Decode(&u)
	if resultA != mongo.ErrNoDocuments {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Account Already Exists"}})
	}
	u.ID = primitive.NewObjectID()
	u.FirstName = body.FirstName
	u.LastName = body.LastName
	u.SetPassword(body.Password)
	u.GenerateRecoveryCode()
	u.CreatedTime = time.Now()
	_, err := dbcontext.UserModel.InsertOne(c.Context(), u)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{err.Error()}})
	}
	ep := emailtemplates.GenerateWelcomeEmailTemplate(*u)
	emailqueue.DispatchEmail(ep)
	return c.JSON(fiber.Map{"data": u})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	body := new(dto.Login)
	c.BodyParser(body)
	validationResult := body.Validate()
	if len(validationResult) != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	u := new(models.User)
	resultA := dbcontext.UserModel.FindOne(c.Context(), bson.D{{Key: "email", Value: body.Email}}).Decode(&u)
	if resultA == mongo.ErrNoDocuments {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"error": []string{"Please check your credentials."}})
	}

	ok := u.VerifyPassword(body.Password)
	if !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"error": []string{"Please check your credentials."}})
	}

	claims := new(jsonwebtoken.AIMClaims)
	claims.ID = string(u.ID.Hex())
	claims.PrincipalID = string(u.ID.Hex())
	claims.PrincipalName = fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	claims.PrincipalEmail = u.Email
	signedStringAccessToken, err := jsonwebtoken.GenerateAccessToken(*claims)
	if err != nil {
		log.Panicln("Cannot Generate Access Token", err)
	}
	signedStringRefreshToken, err := jsonwebtoken.GenerateRefreshToken(*claims)
	if err != nil {
		log.Panicln("Cannot Generate Refresh Token", err)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"access_token":         signedStringAccessToken,
		"access_token_expiry":  time.Now().Add(time.Minute * config.JWT_ACCESS_EXPIRY_MINUTES).UTC(),
		"refresh_token":        signedStringRefreshToken,
		"refresh_token_expiry": time.Now().Add(time.Minute * config.JWT_REFRESH_EXPIRY_MINUTES).UTC()}})
}

func (ac *AuthController) VerifyToken(c *fiber.Ctx) error {
	token := c.Get("authorization")
	if len(token) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Authorization Header Invalid"}})
	}
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Authorization Header Invalid"}})
	}
	claims, err := jsonwebtoken.ValidateToken(tokenParts[1])
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid Token"}})
	}
	return c.JSON(fiber.Map{"data": claims})
}

func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	body := new(dto.RefreshToken)
	c.BodyParser(body)

	validationResult := body.Validate()
	if len(validationResult) != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": validationResult})
	}

	claims, err := jsonwebtoken.ValidateToken(body.RefreshToken)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid Token"}})
	}

	signedStringAccessToken, err := jsonwebtoken.GenerateAccessToken(*claims)
	if err != nil {
		log.Panicln("Cannot Generate Access Token", err)
	}
	signedStringRefreshToken, err := jsonwebtoken.GenerateRefreshToken(*claims)
	if err != nil {
		log.Panicln("Cannot Generate Refresh Token", err)
	}

	return c.JSON(fiber.Map{"data": fiber.Map{
		"access_token":         signedStringAccessToken,
		"access_token_expiry":  time.Now().Add(time.Minute * config.JWT_ACCESS_EXPIRY_MINUTES).UTC(),
		"refresh_token":        signedStringRefreshToken,
		"refresh_token_expiry": time.Now().Add(time.Minute * config.JWT_REFRESH_EXPIRY_MINUTES).UTC()}})
}
