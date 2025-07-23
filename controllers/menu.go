package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// MenuController handles menu CRUD endpoints.
type MenuController struct {
	Repo     *repositories.MenuRepository
	UserRepo *repositories.UserRepository
	SARepo   *repositories.ServiceAccountRepository
}

// NewMenuController creates a new controller using the provided repository.
func NewMenuController(repo *repositories.MenuRepository, userRepo *repositories.UserRepository, saRepo *repositories.ServiceAccountRepository) *MenuController {
	return &MenuController{Repo: repo, UserRepo: userRepo, SARepo: saRepo}
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

	menus, err := ctrl.Repo.GetAll(c.Context(), search)
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

// GetMenusByToken returns menus appropriate for the authenticated entity.
// Service accounts receive only SA menus while regular users receive menus for their unit.
func (ctrl *MenuController) GetMenusByToken(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid token",
			Data:    nil,
		})
	}
	claims := userToken.Claims.(jwt.MapClaims)
	id, _ := claims["id"].(string)

	// try service account first
	sa, err := ctrl.SARepo.FindByID(c.Context(), id)
	if err == nil && sa.Active {
		menus, err := ctrl.Repo.GetSAMenus(c.Context())
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

	// treat as regular user
	user, err := ctrl.UserRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	menus, err := ctrl.Repo.GetUnitMenus(c.Context(), user.UnitID, user.IsAdmin)
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
