package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func DecodedToken(c echo.Context)(map[string]interface{}, error){
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("UnAuthorized")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, fmt.Errorf("UnAuthorized")
	}

	token := tokenParts[1]
	currentUser, err := VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("UnAuthorized")
	}

	return currentUser, nil
}