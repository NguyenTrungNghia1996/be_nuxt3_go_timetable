package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitController struct {
	Repo *repositories.UnitRepository
}

func NewUnitController(repo *repositories.UnitRepository) *UnitController {
	return &UnitController{Repo: repo}
}

func (ctrl *UnitController) Create(c *fiber.Ctx) error {
        var input models.Unit
        if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}
	if _, err := ctrl.Repo.FindBySubDomain(c.Context(), input.SubDomain); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "sub domain exists",
			Data:    nil,
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create unit",
			Data:    err.Error(),
		})
        }
       if input.Active == false {
               input.Active = true
       }
        if err := ctrl.Repo.Create(c.Context(), &input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create unit",
			Data:    err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  "success",
		Message: "Unit created",
		Data:    input,
	})
}

func (ctrl *UnitController) List(c *fiber.Ctx) error {
	if id := c.Query("id"); id != "" {
		unit, err := ctrl.Repo.FindByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Status:  "error",
				Message: "Unit not found",
				Data:    nil,
			})
		}
		return c.JSON(models.APIResponse{
			Status:  "success",
			Message: "Unit retrieved",
			Data:    unit,
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	units, total, err := ctrl.Repo.GetAll(c.Context(), int64(page), int64(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list units",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Units retrieved",
		Data: fiber.Map{
			"items": units,
			"total": total,
		},
	})
}

func (ctrl *UnitController) Get(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing id",
			Data:    nil,
		})
	}
	unit, err := ctrl.Repo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Unit not found",
			Data:    nil,
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Unit retrieved",
		Data:    unit,
	})
}

func (ctrl *UnitController) GetBySubDomain(c *fiber.Ctx) error {
	sub := c.Query("sub_domain")
	if sub == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing sub_domain",
			Data:    nil,
		})
	}
	unit, err := ctrl.Repo.FindBySubDomain(c.Context(), sub)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Unit not found",
			Data:    nil,
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Unit retrieved",
		Data:    unit,
	})
}

func (ctrl *UnitController) Update(c *fiber.Ctx) error {
       var input models.Unit
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

       if u, err := ctrl.Repo.FindBySubDomain(c.Context(), input.SubDomain); err == nil && u.ID.Hex() != id {
               return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
                       Status:  "error",
                       Message: "sub domain exists",
                       Data:    nil,
		})
	} else if err != nil && err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update unit",
			Data:    err.Error(),
		})
	}
       if err := ctrl.Repo.Update(c.Context(), id, &input); err != nil {
               return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
                       Status:  "error",
                       Message: "Failed to update unit",
                       Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Unit updated",
		Data:    input,
	})
}

func (ctrl *UnitController) Delete(c *fiber.Ctx) error {
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
			Message: "Failed to delete unit",
			Data:    err.Error(),
		})
	}
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Unit deleted",
		Data:    nil,
	})
}
