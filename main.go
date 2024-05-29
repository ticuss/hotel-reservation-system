package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/api"
	"github.com/ticuss/hotel-reservation-system/api/middleware"
	"github.com/ticuss/hotel-reservation-system/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userColl = "users"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	now := time.Now()
	fmt.Println(now)
	listenAddr := flag.String("listenAddr", ":5001", "The server listen address")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		// roomHandler  = api.NewRoomHandler(roomStore)
		authHandler = api.NewAuthHandler(userStore)
		roomHandler = api.NewRoomHandler(store)
		app         = fiber.New(config)
		auth        = app.Group("/api")
		apiv1       = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)
	// user handlers
	app.Get("/foo", handleFoo)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user/", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/hotel/", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRooms)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Works good"})
}
