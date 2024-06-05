package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ticuss/hotel-reservation-system/db"
	"github.com/ticuss/hotel-reservation-system/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn string, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.UserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin

	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	fmt.Println("Inserted Used", user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name string, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomsIDS := rooms
	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomsIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, seaside bool, size string, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: seaside,
		Price:   price,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userID primitive.ObjectID, roomID primitive.ObjectID, startDate time.Time, endDate time.Time) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: startDate,
		TillDate: endDate,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
