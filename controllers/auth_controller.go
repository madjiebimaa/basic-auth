package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/madjiebimaa/go-basic-auth/helpers"
	"github.com/madjiebimaa/go-basic-auth/models/domains"
	"github.com/madjiebimaa/go-basic-auth/models/webs"
	"github.com/madjiebimaa/go-basic-auth/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (cr *AuthController) Register(c *fiber.Ctx) error {
	var userRegisterRequest webs.UserRegisterRequest
	err := c.BodyParser(&userRegisterRequest)
	helpers.PanicIfError(err)

	user, err := cr.AuthService.Register(c, userRegisterRequest)
	if err != nil {
		webResponse := webs.WebResponse{
			Code:   fiber.StatusNotAcceptable,
			Status: "NOT OK",
			Data:   err,
		}

		return c.JSON(webResponse)
	}

	userRegisterResponse := webs.UserRegisterResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	webResponse := webs.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   userRegisterResponse,
	}

	return c.JSON(webResponse)
}

func (cr *AuthController) Login(c *fiber.Ctx) error {
	var userLoginRequest webs.UserLoginRequest
	err := c.BodyParser(&userLoginRequest)
	helpers.PanicIfError(err)

	cookie, err := cr.AuthService.Login(c, userLoginRequest)
	if err != nil {
		webResponse := webs.WebResponse{
			Code:   fiber.StatusNotAcceptable,
			Status: "NOT OK",
			Data:   err,
		}

		return c.JSON(webResponse)
	}

	c.Cookie(cookie)

	message := domains.Message{
		Message: "login success",
	}

	webResponse := webs.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   message,
	}

	return c.JSON(webResponse)
}

func (cr *AuthController) User(c *fiber.Ctx) error {
	user, err := cr.AuthService.User(c)
	if err != nil {
		webResponse := webs.WebResponse{
			Code:   fiber.StatusNotAcceptable,
			Status: "NOT OK",
			Data:   err,
		}

		return c.JSON(webResponse)
	}

	userLoginRequest := webs.UserLoginResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	webResponse := webs.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   userLoginRequest,
	}

	return c.JSON(webResponse)

}

func (cr *AuthController) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour), // 1 Day
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	message := domains.Message{
		Message: "logout success",
	}

	webResponse := webs.WebResponse{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   message,
	}

	return c.JSON(webResponse)
}
