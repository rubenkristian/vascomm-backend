package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/pkg"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	userService   *services.UserService
	authGenerator *pkg.AuthToken
	oauth2Config  *configs.ConfigOauth
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func InitializeAuthHandler(userService *services.UserService, authGenerator *pkg.AuthToken, oauth2Config *configs.ConfigOauth) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		authGenerator: authGenerator,
		oauth2Config:  oauth2Config,
	}
}

func (authHandler *AuthHandler) LoginUser(c *fiber.Ctx) error {
	loginInput := new(LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	user, err := authHandler.userService.UserLogin(loginInput.Username, loginInput.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": err.Error(),
			"data":    nil,
		})
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success authorization",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (authHandler *AuthHandler) LoginAdmin(c *fiber.Ctx) error {
	loginInput := new(LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	user, err := authHandler.userService.AdminLogin(loginInput.Username, loginInput.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": err.Error(),
			"data":    nil,
		})
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success authorization",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (authHandler *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshTokenInput := new(RefreshTokenInput)

	if err := c.BodyParser(refreshTokenInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	exp, userId, err := authHandler.authGenerator.ValidateRefresh(refreshTokenInput.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	now := time.Now().Unix()

	if int64(exp) < now {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": "refresh token expired",
			},
		})
	}

	user, err := authHandler.userService.GetUser(uint(userId))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success authorization",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (authHandler *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	key := authHandler.oauth2Config.Key
	url := authHandler.oauth2Config.Config.AuthCodeURL(key, oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func (authHandler *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := authHandler.oauth2Config.Config.Exchange(context.Background(), code)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	client := authHandler.oauth2Config.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	name := userInfo["name"].(string)
	email := userInfo["email"].(string)
	phoneNumber := userInfo["phone"].(string)

	user := authHandler.userService.RegisterFromGoogle(name, email, phoneNumber)

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success authorization",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (authHandler *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	userInput := new(models.User)

	if err := c.BodyParser(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if err := authHandler.userService.CreateUser(userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "Success register user",
		"data":    nil,
	})
}
