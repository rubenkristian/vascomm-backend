package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func InitializeUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (userHandler *UserHandler) GetUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("user_id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	user, err := userHandler.userService.GetUser(uint(userId))

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
		"message": "Something went wrong",
		"data":    user,
	})
}

func (userHandler *UserHandler) GetAllUser(c *fiber.Ctx) error {
	take := c.QueryInt("take", 10)
	skip := c.QueryInt("skip", 0)
	search := c.Query("search", "")

	users, err := userHandler.userService.GetAllUser(take, skip, search)

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
		"message": "Success get all users",
		"data":    users,
	})
}

func (userHandler *UserHandler) CreateUser(c *fiber.Ctx) error {
	input := new(models.User)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if err := userHandler.userService.CreateUser(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusInternalServerError,
		"message": "Something went wrong",
		"data":    input,
	})
}

func (userHandler *UserHandler) UpdateUser(c *fiber.Ctx) error {
	input := new(models.User)
	userId, err := c.ParamsInt("user_id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	user, err := userHandler.userService.UpdateUser(uint(userId), input)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "Success update user",
		"data":    user,
	})
}

func (userHandler *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("user_id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if err := userHandler.userService.DeleteUser(uint(userId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "Success delete user",
		"data":    nil,
	})
}
