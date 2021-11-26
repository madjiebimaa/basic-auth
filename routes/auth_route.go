package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/madjiebimaa/go-basic-auth/controllers"
	"github.com/madjiebimaa/go-basic-auth/services"
	"gorm.io/gorm"
)

func NewAuthRoute(db *gorm.DB, app *fiber.App, validate *validator.Validate) {

	authService := services.NewAuthService(db, validate)
	authController := controllers.NewAuthController(authService)

	app.Get("/api/users", authController.User)
	app.Post("/api/register", authController.Register)
	app.Post("/api/login", authController.Login)
	app.Post("/api/logout", authController.Logout)
}
