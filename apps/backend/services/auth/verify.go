package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func IsCurrentUser(ctx *fiber.Ctx, userId uuid.UUID) bool {
	user := ctx.Locals("user").(*JwtClaims)
	return user.UserId == userId
}
