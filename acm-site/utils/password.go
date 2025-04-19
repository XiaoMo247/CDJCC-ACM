package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

func CheckPassword(pw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
	return err == nil
}
