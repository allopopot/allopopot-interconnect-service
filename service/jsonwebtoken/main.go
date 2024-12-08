package jsonwebtoken

import (
	"allopopot-interconnect-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AIMClaims struct {
	PrincipalID    int
	PrincipalName  string
	PrincipalEmail string
	jwt.RegisteredClaims
}

func GenerateToken(claims AIMClaims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * config.JWT_EXPIRY_MINUTES))
	claims.Issuer = config.JWT_ISSUER
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(config.JWT_SECRET))
	return signedString, err
}

func ValidateToken(tokenString string) (*AIMClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AIMClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})
	return token.Claims.(*AIMClaims), err
}
