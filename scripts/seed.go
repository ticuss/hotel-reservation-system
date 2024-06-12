package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ticuss/hotel-reservation-system/api"
	"github.com/ticuss/hotel-reservation-system/db"
	"github.com/ticuss/hotel-reservation-system/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Clinet       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func main() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	mongoDbName := os.Getenv(db.MongoDBNameEnvName)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)

	userStore = db.NewMongoUserStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	admin := fixtures.AddUser(store, "kek", "kaka", true)
	fmt.Printf("User: %v\n", api.CreateTokenFromUser(user))
	fmt.Printf("User: %v\n", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "TopKeke", "Asnieres", 1, nil)
	room := fixtures.AddRoom(store, true, "small", 99.9, hotel.ID)
	_ = fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 1))
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Random Hotel name: %v", i)
		location := fmt.Sprintf("Random Hotel location: %v", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}

func init() {
}
