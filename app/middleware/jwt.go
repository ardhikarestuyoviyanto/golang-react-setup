package middleware

import (
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
