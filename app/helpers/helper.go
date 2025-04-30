package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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

func Hash(plainText string)(string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyHash(hashed, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plainText))
}

func EncryptString(plainText string)(string, error){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error load env file")
	}
	
	var appKey = os.Getenv("APP_KEY")

	if len(appKey) != 32{
		return "", errors.New("APP_KEY harus 32 digits")
	}

	var secretKey = []byte(appKey)

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
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error load env file")
	}
	
	var appKey = os.Getenv("APP_KEY")

	if len(appKey) != 32{
		return "", errors.New("APP_KEY harus 32 digits")
	}

	var secretKey = []byte(appKey)

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

func UploadFile(dstDir string, file *multipart.FileHeader)(string, error){
	ext := filepath.Ext(file.Filename)
	randomFileName := uuid.New().String() + ext
	dstPath := filepath.Join(dstDir, randomFileName)

	src, err := file.Open()
	if err != nil {
		return "", errors.New("Gagal Membuka File")
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New("Gagal Menyimpan File File")
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", errors.New("Gagal Menyalin File File")
	}

	return randomFileName, nil
}

func ValidateExtFile(extAllowed []string, file *multipart.FileHeader, maxSizeInMb int) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Validasi ekstensi
	allowed := false
	for _, allowedExt := range extAllowed {
		if ext == strings.ToLower(allowedExt) {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("ekstensi file harus salah satu dari: %v", extAllowed)
	}

	// Validasi ukuran file
	if file.Size > int64(maxSizeInMb)*1024*1024 {
		return fmt.Errorf("ukuran file maksimal %d MB", maxSizeInMb)
	}

	return nil
}
// func SetDefault(value string, defaultValue string) string {
//     if value == "" {
//         return defaultValue
//     }
//     return value
// }