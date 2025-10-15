package auth

import (
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *AuthService) IsCurrentUser(ctx *fiber.Ctx, userId uuid.UUID) (bool, error) {
	user, err := s.RequireUserContext(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user context")
	}
	return user.UserId == userId, nil
}
