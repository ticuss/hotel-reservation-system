package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Foo",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "James Foo"})
}
