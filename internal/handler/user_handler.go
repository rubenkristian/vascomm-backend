package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/utils"
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
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	user, err := userHandler.userService.GetUser(uint(userId))

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success get user", user)(c)
}

func (userHandler *UserHandler) GetAllUser(c *fiber.Ctx) error {
	take := c.QueryInt("take", 10)
	skip := c.QueryInt("skip", 0)
	search := c.Query("search", "")
	sort := c.Query("sort", "asc")
	sortBy := c.Query("sortBy", "id")

	users, err := userHandler.userService.GetAllUser(take, skip, search, sort, sortBy)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success get all users", users)(c)
}

func (userHandler *UserHandler) CreateUser(c *fiber.Ctx) error {
	input := new(models.User)

	if err := c.BodyParser(input); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	if err := userHandler.userService.CreateUser(input); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusCreated, "Success create user", input)(c)
}

func (userHandler *UserHandler) UpdateUser(c *fiber.Ctx) error {
	input := new(models.User)
	userId, err := c.ParamsInt("user_id")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	if err := c.BodyParser(input); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	user, err := userHandler.userService.UpdateUser(uint(userId), input)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusCreated, "Success update user", user)(c)
}

func (userHandler *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("user_id")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	if err := userHandler.userService.DeleteUser(uint(userId)); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusCreated, "Success delete user", nil)(c)
}
