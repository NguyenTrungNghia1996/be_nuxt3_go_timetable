package controllers

// This file contains authentication handlers used by the API.

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication requests.
type AuthController struct {
	UserRepo *repositories.UserRepository
	UnitRepo *repositories.UnitRepository
}

// NewAuthController creates a controller with the provided user repository.
func NewAuthController(userRepo *repositories.UserRepository, unitRepo *repositories.UnitRepository) *AuthController {
	return &AuthController{UserRepo: userRepo, UnitRepo: unitRepo}
}

// Login authenticates a user and returns a signed JWT on success.
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UnitID   string `json:"unit_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}

	if input.UnitID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing unit_id",
			Data:    nil,
		})
	}

	user, err := ctrl.UserRepo.FindByUsername(c.Context(), input.Username)
	if err != nil || user.UnitID.Hex() != input.UnitID || !user.Active || !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid credentials",
			Data:    nil,
		})
	}

	unit, err := ctrl.UnitRepo.FindByID(c.Context(), input.UnitID)
	if err != nil || !unit.Active {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid credentials",
			Data:    nil,
		})
	}

	token, _ := utils.GenerateJWT(user.ID.Hex())
	u := fiber.Map{
		"id":         user.ID.Hex(),
		"username":   user.Username,
		"name":       user.Name,
		"url_avatar": user.UrlAvatar,
		"unit_id":    user.UnitID.Hex(),
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Login successful",
		Data: fiber.Map{
			"token": token,
			"user":  u,
		},
	})
}
