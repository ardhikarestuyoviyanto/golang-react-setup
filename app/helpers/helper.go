package helpers

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Regex email sederhana
func ValidateEmail(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Email tidak valid")
	}
	return nil
}

// Password minimal 6 karakter, bisa kamu sesuaikan
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("Password minimal 6 karakter")
	}
	return nil
}

func HashString (plainText string)(string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHash(hashed, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plainText))
}