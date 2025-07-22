package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuController struct {
	Repo *repositories.MenuRepository
}

func NewMenuController(repo *repositories.MenuRepository) *MenuController {
	return &MenuController{Repo: repo}
}

func (ctrl *MenuController) CreateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid data",
			Data:    nil,
		})
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	unitIDHex, _ := claims["unit_id"].(string)
	menu.UnitID, _ = primitive.ObjectIDFromHex(unitIDHex)

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

func (ctrl *MenuController) GetMenus(c *fiber.Ctx) error {
	search := c.Query("search")
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	unitIDHex, _ := claims["unit_id"].(string)
	uid, _ := primitive.ObjectIDFromHex(unitIDHex)

	menus, err := ctrl.Repo.GetAll(c.Context(), uid, search)
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

func (ctrl *MenuController) DeleteMenu(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	unitIDHex, _ := claims["unit_id"].(string)
	uid, _ := primitive.ObjectIDFromHex(unitIDHex)

	if err := ctrl.Repo.DeleteByID(c.Context(), uid, id); err != nil {
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

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	unitIDHex, _ := claims["unit_id"].(string)
	uid, _ := primitive.ObjectIDFromHex(unitIDHex)

	if err := ctrl.Repo.UpdateByID(c.Context(), uid, menu.ID.Hex(), &menu); err != nil {
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
