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
	"go.mongodb.org/mongo-driver/mongo/options"
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

	accessTokenExpiry := time.Now().Add(time.Minute * config.JWT_ACCESS_EXPIRY_MINUTES).UTC()
	refreshTokenExpiry := time.Now().Add(time.Minute * config.JWT_REFRESH_EXPIRY_MINUTES).UTC()

	t := new(models.Token)
	t.CreatedTime = time.Now()
	t.Type = "refresh"
	t.Token = signedStringRefreshToken
	t.UserID = u.ID
	t.ExpiryTime = refreshTokenExpiry
	dbcontext.TokenModel.ReplaceOne(c.Context(), bson.D{{Key: "user_id", Value: u.ID}}, t, options.Replace().SetUpsert(true))

	return c.JSON(fiber.Map{"data": fiber.Map{
		"access_token":         signedStringAccessToken,
		"access_token_expiry":  accessTokenExpiry,
		"refresh_token":        signedStringRefreshToken,
		"refresh_token_expiry": refreshTokenExpiry,
	}})
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

	t := new(models.Token)
	claimIdInObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid User ID"}})
	}

	findToken := dbcontext.TokenModel.FindOne(c.Context(), bson.D{{Key: "user_id", Value: claimIdInObjectId}, {Key: "expiry_time", Value: bson.D{{Key: "$gte", Value: time.Now()}}}}).Decode(&t)
	if findToken != nil {
		if findToken == mongo.ErrNoDocuments {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": []string{"Token Expired"}})
		}
	}

	if t.Token != body.RefreshToken {
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

	accessTokenExpiry := time.Now().Add(time.Minute * config.JWT_ACCESS_EXPIRY_MINUTES).UTC()
	refreshTokenExpiry := time.Now().Add(time.Minute * config.JWT_REFRESH_EXPIRY_MINUTES).UTC()

	t.CreatedTime = time.Now()
	t.Token = signedStringRefreshToken
	t.ExpiryTime = refreshTokenExpiry
	dbcontext.TokenModel.ReplaceOne(c.Context(), bson.D{{Key: "user_id", Value: t.UserID}}, t, options.Replace().SetUpsert(true))

	return c.JSON(fiber.Map{"data": fiber.Map{
		"access_token":         signedStringAccessToken,
		"access_token_expiry":  accessTokenExpiry,
		"refresh_token":        signedStringRefreshToken,
		"refresh_token_expiry": refreshTokenExpiry,
	}})
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
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

	claimIdInObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid User ID"}})
	}

	dbcontext.TokenModel.DeleteMany(c.Context(), bson.M{"user_id": claimIdInObjectId})

	c.ClearCookie("access_token", "refresh_token")
	return c.JSON(fiber.Map{"data": true})
}

func (ac *AuthController) DeleteAccount(c *fiber.Ctx) error {

	access_token := ""

	tokenFromCookie := c.Cookies("access_token")
	tokenFromAuthHeader := strings.Split(c.Get("authorization"), " ")

	if len(tokenFromCookie) > 0 {
		access_token = tokenFromCookie
	}

	if len(tokenFromAuthHeader) == 2 {
		access_token = tokenFromAuthHeader[1]
	}

	claims, err := jsonwebtoken.ValidateToken(access_token)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid Token"}})
	}

	claimIdInObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": []string{"Invalid User ID"}})
	}

	dbcontext.UserModel.DeleteOne(c.Context(), bson.M{"_id": claimIdInObjectId})
	dbcontext.TokenModel.DeleteMany(c.Context(), bson.M{"user_id": claimIdInObjectId})

	c.ClearCookie("access_token", "refresh_token")
	return c.JSON(fiber.Map{"data": true})
}
