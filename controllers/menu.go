package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"strings"
)

// MenuController handles menu CRUD endpoints.
type MenuController struct {
	Repo *repositories.MenuRepository
}

// NewMenuController creates a new controller using the provided repository.
func NewMenuController(repo *repositories.MenuRepository) *MenuController {
	return &MenuController{Repo: repo}
}

// CreateMenu parses and persists a new menu item.
func (ctrl *MenuController) CreateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid data",
			Data:    nil,
		})
	}

	if err := ctrl.Repo.Create(c.Context(), &menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Unable to create menu",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Created menu successfully",
		Data:    menu.ToResponse(),
	})
}

// GetMenus returns a list of menus optionally filtered by search keyword.
func (ctrl *MenuController) GetMenus(c *fiber.Ctx) error {
	search := c.Query("search")

	// Determine SA filter
	var isSA *bool
	if q := c.Query("is_sa"); q != "" {
		if strings.ToLower(q) == "true" {
			t := true
			isSA = &t
		} else {
			f := false
			isSA = &f
		}
	}

	menus, err := ctrl.Repo.GetAll(c.Context(), search, isSA)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Cannot get menu list",
			Data:    nil,
		})
	}

	resp := make([]models.MenuResponse, len(menus))
	for i, m := range menus {
		resp[i] = m.ToResponse()
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Get menu list successfully",
		Data:    resp,
	})
}

// DeleteMenu removes a menu item by id.
func (ctrl *MenuController) DeleteMenu(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}

	if err := ctrl.Repo.DeleteByID(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Deleted menu successfully",
		Data:    nil,
	})
}

// UpdateMenu updates an existing menu item.
func (ctrl *MenuController) UpdateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid data",
			Data:    nil,
		})
	}

	if menu.ID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}

	if err := ctrl.Repo.UpdateByID(c.Context(), menu.ID.Hex(), &menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Updated menu successfully",
		Data:    menu.ToResponse(),
	})
}
