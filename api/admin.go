package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if !user.IsAdmin {
		return ErrUnauthorized()
	}
	return c.Next()
}
