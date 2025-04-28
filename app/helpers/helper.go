package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var appKey = os.Getenv("APP_KEY")
var secretKey = []byte(appKey)

func ValidateEmail(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Email tidak valid")
	}
	return nil
}

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

func EncryptString(plainText string)(string, error){
	block , err := aes.NewCipher(secretKey)
	if err != nil{
		return "", err
	}

	chipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := chipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil{
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(chipherText[aes.BlockSize:], []byte(plainText))
	return base64.URLEncoding.EncodeToString(chipherText), nil
}

func DecryptString(encryptedText string)(string, error){
	cipherText, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil{
		return "",err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil{
		return "",err
	}

	if len(cipherText) < aes.BlockSize{
		return "", errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// func SetDefault(value string, defaultValue string) string {
//     if value == "" {
//         return defaultValue
//     }
//     return value
// }