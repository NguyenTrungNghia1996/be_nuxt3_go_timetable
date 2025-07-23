package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

type ServiceAccountController struct {
	Repo *repositories.ServiceAccountRepository
}

func NewServiceAccountController(repo *repositories.ServiceAccountRepository) *ServiceAccountController {
	return &ServiceAccountController{Repo: repo}
}

func (ctrl *ServiceAccountController) Create(c *fiber.Ctx) error {
	var input models.ServiceAccount
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
			Message: "Failed to create service account",
			Data:    err.Error(),
		})
	} else if exists {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "username exists",
			Data:    nil,
		})
	}
	if !input.Active {
		input.Active = true
	}
	if input.Password != "" {
		hashed, _ := utils.HashPassword(input.Password)
		input.Password = hashed
	}
	if err := ctrl.Repo.Create(c.Context(), &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create service account",
			Data:    err.Error(),
		})
	}
	input.Password = ""
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  "success",
		Message: "Service account created",
		Data:    input,
	})
}

func (ctrl *ServiceAccountController) List(c *fiber.Ctx) error {
	if id := c.Query("id"); id != "" {
		item, err := ctrl.Repo.FindByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "Service account not found",
				Data:    nil,
			})
		}
		return c.JSON(models.APIResponse{
			Status:  "success",
			Message: "Service account retrieved",
			Data:    item,
		})
	}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	items, total, err := ctrl.Repo.GetAll(c.Context(), int64(page), int64(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list service accounts",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Service accounts retrieved",
		Data: fiber.Map{
			"items": items,
			"total": total,
		},
	})
}

func (ctrl *ServiceAccountController) Update(c *fiber.Ctx) error {
	var input models.ServiceAccount
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
	id := input.ID.Hex()
	if input.Password != "" {
		hashed, _ := utils.HashPassword(input.Password)
		input.Password = hashed
	}
	if err := ctrl.Repo.Update(c.Context(), id, &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update service account",
			Data:    err.Error(),
		})
	}
	input.Password = ""
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Service account updated",
		Data:    input,
	})
}

func (ctrl *ServiceAccountController) Delete(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}
	if err := ctrl.Repo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete service account",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Service account deleted",
		Data:    nil,
	})
}
