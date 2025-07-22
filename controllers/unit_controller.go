package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UnitController struct {
	Repo *repositories.UnitRepository
}

func NewUnitController(repo *repositories.UnitRepository) *UnitController {
	return &UnitController{Repo: repo}
}

func (ctrl *UnitController) CreateUnit(c *fiber.Ctx) error {
	var unit models.Unit
	if err := c.BodyParser(&unit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid data",
			Data:    nil,
		})
	}
	if err := ctrl.Repo.Create(c.Context(), &unit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Unable to create unit",
			Data:    nil,
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Created unit successfully",
		Data:    unit.ToResponse(),
	})
}

func (ctrl *UnitController) GetUnits(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	units, total, err := ctrl.Repo.GetAll(c.Context(), search, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Cannot get unit list",
			Data:    nil,
		})
	}
	resp := make([]models.UnitResponse, len(units))
	for i, u := range units {
		resp[i] = u.ToResponse()
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Get unit list successfully",
		Data: fiber.Map{
			"items": resp,
			"total": total,
		},
	})
}

func (ctrl *UnitController) UpdateUnit(c *fiber.Ctx) error {
	var unit models.Unit
	if err := c.BodyParser(&unit); err != nil || unit.ID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid data",
			Data:    nil,
		})
	}
	if err := ctrl.Repo.UpdateByID(c.Context(), unit.ID.Hex(), &unit); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "unit not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Updated unit successfully",
		Data:    unit.ToResponse(),
	})
}

func (ctrl *UnitController) DeleteUnit(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}
	if err := ctrl.Repo.DeleteByID(c.Context(), id); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "unit not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Deleted unit successfully",
		Data:    nil,
	})
}
