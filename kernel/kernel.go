package kernel

import (
	"go-auth/app/models"
	"go-auth/app/routes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}



func StartApplication(e *echo.Echo) {
	// Load Env
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error load env file")
	}

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
	// Register views
	basePath, _ := os.Getwd()
	t := &Template{
		templates: template.Must(template.ParseGlob(filepath.Join(basePath, "app/views/**/*.html"))),
	}
	e.Renderer = t

	// Register static file untuk frontend
	e.Static("/assets", "app/views/js/dist/assets")
	e.Static("/", "app/views/js/public")

	// Debug App
	if appEnv == "development" {
		e.Debug = true
	} else {
		e.Debug = false
	}	// Start Application
	e.Logger.Fatal(e.Start(port))
}

