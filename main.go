package main

import (
	"go-auth/kernel"

	"github.com/labstack/echo/v4"
)

func main() {
	// Start Application
	e := echo.New()
	kernel.StartApplication(e)
}