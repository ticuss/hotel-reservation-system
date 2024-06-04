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
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	app.Get("/foo", handleFoo)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room handlers
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Post("/room", roomHandler.HandleGetRooms)

	// booking handlers
	apiv1.Post("/booking", roomHandler.HandleBookRoom)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Works good"})
}
