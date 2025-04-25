package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var studentSecret = []byte("your-student-secret-key") // 可以和管理员的 secret 区分开

type StudentClaims struct {
	StudentID     uint   `json:"student_id"`
	StudentNumber string `json:"student_number"`
	jwt.RegisteredClaims
}

// 生成 Token
func GenerateStudentToken(studentID uint, studentNumber string) (string, error) {
	claims := StudentClaims{
		StudentID:     studentID,
		StudentNumber: studentNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(studentSecret)
}

// 解析 Token
func ParseStudentToken(tokenStr string) (*StudentClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &StudentClaims{}, func(token *jwt.Token) (interface{}, error) {
		return studentSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*StudentClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
