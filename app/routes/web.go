package routes

import (
	"go-auth/app/handler/auth"
	"go-auth/app/handler/auth/captcha_"
	"go-auth/app/handler/task"
	_midlleware2 "go-auth/app/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func RoutesWeb(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	
	e.GET("/api/v1/captcha/generateCaptcha", captcha_.GenerateCaptcha)
	e.GET("/api/v1/captcha/:captchaId/get", captcha_.GetCaptcha)
	e.POST("/api/v1/signUp", auth.SignUpHandler(db))
	e.POST("/api/v1/signIn", auth.SignInHandler(db))
	// With Middleware Auth
	e.GET("/api/v1/tasks", task.GetAllHandler(db), _midlleware2.AuthMidleware)
	e.POST("/api/v1/tasks", task.StoreHandler(db), _midlleware2.AuthMidleware)
	e.GET("/api/v1/tasks/:taskId", task.GetHandler(db), _midlleware2.AuthMidleware)
	e.DELETE("/api/v1/tasks/:taskId", task.DestroyHandler(db), _midlleware2.AuthMidleware)
	e.PUT("/api/v1/tasks/:taskId", task.UpdateHandler(db), _midlleware2.AuthMidleware)

}