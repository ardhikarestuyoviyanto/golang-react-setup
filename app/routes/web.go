package routes

import (
	"go-auth/app/handler/auth"
	"go-auth/app/handler/auth/captcha_"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func RoutesWeb(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.CORS()) // <- ini penting banget
	e.GET("/", auth.LoginView)
	e.GET("/api/v1/captcha/generateCaptcha", captcha_.GenerateCaptcha)
	e.GET("/api/v1/captcha/:captchaId/get", captcha_.GetCaptcha)
	e.POST("/api/v1/signUp", auth.SignUpHandler(db))
	e.POST("/api/v1/signIn", auth.SignInHandler(db))


}