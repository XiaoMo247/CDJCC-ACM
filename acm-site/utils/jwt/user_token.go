package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var userJwtKey = []byte("user_secret_key")

type UserClaims struct {
	UserID        uint
	StudentNumber string
	jwt.RegisteredClaims
}

func GenerateUserToken(userID uint, studentNumber string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &UserClaims{
		UserID:        userID,
		StudentNumber: studentNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(userJwtKey)
}

func ParseUserToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return userJwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*UserClaims), nil
}
