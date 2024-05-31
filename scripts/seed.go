package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ticuss/hotel-reservation-system/db"
	"github.com/ticuss/hotel-reservation-system/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Clinet     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname string, lname string, email string, isAdmin bool, password string) {
	user, err := types.NewUserFromParams(types.UserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	user, err = userStore.InsertUser(context.TODO(), user)
	fmt.Println("Inserted Used", user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},

		{
			Size:  "small",
			Price: 1999.9,
		},

		{
			Size:  "kingsize",
			Price: 129999,
		},
	}

	insetedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seeding the database...")
	for _, room := range rooms {
		room.HotelID = hotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted hotel:", insetedHotel)
		fmt.Println("Inserted hotel:", insertedRoom)
	}
}

func main() {
	seedHotel("Bellucia", "Chisinau", 4)
	seedHotel("TopKeke", "Asnieres", 1)
	seedUser("james", "foo", "kek@kek.com", false, "supersecret")
	seedUser("admin", "admin", "admin@admin.com", true, "supersecret")
}

func init() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)

	userStore = db.NewMongoUserStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
