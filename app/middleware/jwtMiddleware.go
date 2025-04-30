package middleware

import (
	"errors"
	"go-auth/app/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var appKey = os.Getenv("APP_KEY")
var jwtKey = []byte(appKey)

func GenerateJwt(user models.Users) (string, error){
	claims:= jwt.MapClaims{
		"user": map[string]interface{}{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
		},
		"exp": time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyToken(tokenString string)(jwt.MapClaims, error){
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}