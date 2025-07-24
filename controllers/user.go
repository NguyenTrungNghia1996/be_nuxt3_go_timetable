package controllers

import (
	"errors"
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

// UserController handles user CRUD operations.
type UserController struct {
	Repo *repositories.UserRepository
}

// NewUserController creates a new controller instance.
func NewUserController(repo *repositories.UserRepository) *UserController {
	return &UserController{Repo: repo}
}

// Create adds a new user under the admin's unit.
func (ctrl *UserController) Create(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}
	if exists, err := ctrl.Repo.IsUsernameExists(c.Context(), input.Username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create user",
			Data:    err.Error(),
		})
	} else if exists {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "username exists",
			Data:    nil,
		})
	}
	admin := c.Locals("admin_user").(*models.User)
	input.UnitID = admin.UnitID
	if input.Password != "" {
		hashed, _ := utils.HashPassword(input.Password)
		input.Password = hashed
	}
	if err := ctrl.Repo.Create(c.Context(), &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create user",
			Data:    err.Error(),
		})
	}
	input.Password = ""
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  "success",
		Message: "User created",
		Data:    input,
	})
}

// List returns users belonging to the admin's unit.
func (ctrl *UserController) List(c *fiber.Ctx) error {
	admin := c.Locals("admin_user").(*models.User)
	if id := c.Query("id"); id != "" {
		user, err := ctrl.Repo.FindByID(c.Context(), id)
		if err != nil || user.UnitID != admin.UnitID {
			if errors.Is(err, repositories.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
					Status:  "error",
					Message: "User not found",
					Data:    nil,
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "User not found",
				Data:    nil,
			})
		}
		user.Password = ""
		return c.JSON(models.APIResponse{
			Status:  "success",
			Message: "User retrieved",
			Data:    user,
		})
	}
	search := c.Query("search")
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 0)
	users, total, err := ctrl.Repo.GetAllByUnit(c.Context(), admin.UnitID, search, int64(page), int64(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list users",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Users retrieved",
		Data: fiber.Map{
			"items": users,
			"total": total,
		},
	})
}

// Update modifies an existing user within the admin's unit.
func (ctrl *UserController) Update(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}
	if input.ID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}
	admin := c.Locals("admin_user").(*models.User)
	id := input.ID.Hex()
	user, err := ctrl.Repo.FindByID(c.Context(), id)
	if err != nil || user.UnitID != admin.UnitID {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "User not found",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "User not found",
			Data:    nil,
		})
	}
	if input.Password != "" {
		hashed, _ := utils.HashPassword(input.Password)
		input.Password = hashed
	}
	if err := ctrl.Repo.Update(c.Context(), id, &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update user",
			Data:    err.Error(),
		})
	}
	input.Password = ""
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "User updated",
		Data:    input,
	})
}

// Delete removes a user from the admin's unit.
func (ctrl *UserController) Delete(c *fiber.Ctx) error {
	admin := c.Locals("admin_user").(*models.User)
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}
	user, err := ctrl.Repo.FindByID(c.Context(), id)
	if err != nil || user.UnitID != admin.UnitID {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "User not found",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "User not found",
			Data:    nil,
		})
	}
	if err := ctrl.Repo.Delete(c.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "User not found",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete user",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "User deleted",
		Data:    nil,
	})
}
