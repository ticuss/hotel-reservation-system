package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ticuss/hotel-reservation-system/db"
	"github.com/ticuss/hotel-reservation-system/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)

	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Hotel California",
		Location: "California",
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

	insetedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	for _, room := range rooms {
		room.HotelID = hotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Seeding the database...")
		fmt.Println("Inserted hotel:", insetedHotel)
		fmt.Println("Inserted hotel:", insertedRoom)
	}
}
