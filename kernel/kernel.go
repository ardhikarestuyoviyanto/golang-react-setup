package kernel

import (
	"go-auth/app/models"
	"go-auth/app/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func StartApplication(e *echo.Echo) {
	// Load Env
	godotenv.Load()
	// Get Variable
	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Init DB
	configDb := map[string]interface{}{
		"dbHost":    dbHost,
		"dbUser":     dbUser,
		"dbPassword": dbPassword,
		"dbName":     dbName,
		"dbPort":     dbPort,
	}
	db, err := models.InitDb(configDb)

	if err != nil{
		log.Fatal(err)
	}

	// Register routes
	routes.RoutesWeb(e, db)

	// Register static file untuk frontend
	e.Static("/assets", "app/views/js/dist/assets")
	e.Static("/storage", "app/views/storage")
	e.File("/", "app/views/js/dist/index.html")
	e.GET("/*", func(c echo.Context) error {
		return c.File("app/views/js/dist/index.html")
	})

	// Debug App
	if appEnv == "development" {
		e.Debug = true
	} else {
		e.Debug = false
	}	
	// Start Application
	e.Logger.Fatal(e.Start(":" + port))}

