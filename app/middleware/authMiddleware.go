package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMidleware(next echo.HandlerFunc)echo.HandlerFunc{
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": "UnAuthorized"})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": "UnAuthorized"})
		}

		token := tokenParts[1]
		currentUser, err := VerifyToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": "UnAuthorized"})
		}

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": err.Error()})
		}
		c.Set("user", currentUser["user"])
		return next(c)
	}
}