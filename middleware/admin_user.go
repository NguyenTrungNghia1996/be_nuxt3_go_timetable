package middleware

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AdminUser verifies the JWT belongs to an active admin user.
func AdminUser(repo *repositories.UserRepository) fiber.Handler {
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
		user, err := repo.FindByID(c.Context(), id)
		if err != nil || !user.Active || !user.IsAdmin {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Status:  "error",
				Message: "Unauthorized",
				Data:    nil,
			})
		}
		// attach user to context
		c.Locals("admin_user", user)
		return c.Next()
	}
}
