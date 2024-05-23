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
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},

		{
			Type:      types.DeluxeRoomType,
			BasePrice: 1999.9,
		},

		{
			Type:      types.SeaSideRoomType,
			BasePrice: 122.9,
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

	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
