package middleware

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ServiceAccount verifies the JWT belongs to an active service account.
func ServiceAccount(repo *repositories.ServiceAccountRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		sa, err := repo.FindByID(c.Context(), id)
		if err != nil || !sa.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Status:  "error",
				Message: "Unauthorized",
				Data:    nil,
			})
		}
		// attach service account to context for handlers
		c.Locals("service_account", sa)
		return c.Next()
	}
}
