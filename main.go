package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "The server listen address")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("api/v1")

	app.Get("/foo", handleFoo)
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Works good"})
}
