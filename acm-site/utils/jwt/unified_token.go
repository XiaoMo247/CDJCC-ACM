package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 统一的 JWT 密钥
var unifiedJwtKey = []byte("your_secret_key")

// 统一的 Token Claims 结构
type UnifiedClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" | "student" | "member"
	jwt.RegisteredClaims
}

// GenerateUnifiedToken 生成统一的 JWT token
// role: "admin" | "student" | "member"
func GenerateUnifiedToken(userID uint, username string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时过期

	claims := &UnifiedClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "acm-site",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(unifiedJwtKey)
}

// ParseUnifiedToken 解析并验证统一的 token
func ParseUnifiedToken(tokenString string) (*UnifiedClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UnifiedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return unifiedJwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UnifiedClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
