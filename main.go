package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/madjiebimaa/go-basic-auth/apps"
	"github.com/madjiebimaa/go-basic-auth/helpers"
	"github.com/madjiebimaa/go-basic-auth/routes"
)

func main() {
	db := apps.ConnectDB()
	app := fiber.New()
	validate := validator.New()

	routes.NewAuthRoute(db, app, validate)

	err := app.Listen(":3000")
	helpers.PanicIfError(err)
}
