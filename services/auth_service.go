package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/madjiebimaa/go-basic-auth/helpers"
	"github.com/madjiebimaa/go-basic-auth/models/domains"
	"github.com/madjiebimaa/go-basic-auth/models/webs"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewAuthService(db *gorm.DB, validate *validator.Validate) *AuthService {
	return &AuthService{
		DB:       db,
		Validate: validate,
	}
}

func (s *AuthService) Register(c *fiber.Ctx, userRegisterRequest webs.UserRegisterRequest) (*domains.User, error) {
	err := s.Validate.Struct(userRegisterRequest)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil, errors.New("use the right email format and strong password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegisterRequest.Password), 12)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil, errors.New("password can'nt hashed")
	}

	user := domains.User{
		Name:     userRegisterRequest.Name,
		Email:    userRegisterRequest.Email,
		Password: string(hashedPassword),
	}

	s.DB.Create(&user)

	return &user, nil
}

func (s *AuthService) Login(c *fiber.Ctx, userLoginRequest webs.UserLoginRequest) (*fiber.Cookie, error) {
	err := s.Validate.Struct(userLoginRequest)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil, errors.New("use the right email format and strong password")
	}

	var user domains.User
	s.DB.Where("email = ?", userLoginRequest.Email).First(&user)
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLoginRequest.Password))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil, errors.New("incorrect password")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(helpers.SECRET_KEY))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil, errors.New("could not login")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	return &cookie, nil
}

func (s *AuthService) User(c *fiber.Ctx) (*domains.User, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(helpers.SECRET_KEY), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return nil, errors.New("unauthenticated")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user domains.User
	s.DB.Where("id = ?", claims.Issuer).First(&user)

	return &user, nil
}
