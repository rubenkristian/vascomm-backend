package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/pkg"
)

type AuthMiddleware struct {
	userService   *services.UserService
	authGenerator *pkg.AuthToken
}

func InitializeAuthMiddleware(userService *services.UserService, authGenerator *pkg.AuthToken) *AuthMiddleware {
	return &AuthMiddleware{
		userService:   userService,
		authGenerator: authGenerator,
	}
}

func (authMiddleware *AuthMiddleware) CheckAuthorization(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Invalid token format",
			"data":    nil,
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	exp, userId, err := authMiddleware.authGenerator.ValidateToken(tokenString)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	now := time.Now().Unix()

	if int64(exp) < now {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": "Expired token",
			},
		})
	}

	user, err := authMiddleware.userService.GetUser(uint(userId))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": "User not found",
			},
		})
	}

	c.Locals("user", user)

	return c.Next()
}

func (authMiddleware *AuthMiddleware) CheckRole(requiredRole []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user").(*models.User)

		if userData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "Unauthorized",
				"data": fiber.Map{
					"error": "User not found",
				},
			})
		}

		for _, role := range requiredRole {
			if role == userData.Role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Unauthorized",
			"data": fiber.Map{
				"error": "User not authorized to this route",
			},
		})
	}
}
