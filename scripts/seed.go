package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ticuss/hotel-reservation-system/api"
	"github.com/ticuss/hotel-reservation-system/db"
	"github.com/ticuss/hotel-reservation-system/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func seedUser(fname string, lname string, email string, isAdmin bool, password string) *types.User {
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
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return user
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insetedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insetedHotel
}

func seedRoom(hotelID primitive.ObjectID, size string, seaside bool, price float64) *types.Room {
	room := types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: seaside,
		Price:   price,
	}
	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func main() {
	u1 := seedUser("james", "foo", "kek@kek.com", false, "supersecret")
	u2 := seedUser("admin", "admin", "admin@admin.com", true, "supersecret")
	h1 := seedHotel("TopKeke", "Asnieres", 1)
	h2 := seedHotel("Bellucia", "Chisinau", 4)
	h3 := seedHotel("TopKeke", "Asnieres", 1)
	r1 := seedRoom(h1.ID, "small", true, 99.9)
	r2 := seedRoom(h2.ID, "small", true, 199.9)
	r3 := seedRoom(h2.ID, "kingsize", false, 129999)
	r4 := seedRoom(h3.ID, "kingsize", false, 129999)
	seedBooking(u1.ID, r1.ID, time.Now(), time.Now().AddDate(0, 0, 1))
	seedBooking(u1.ID, r2.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	seedBooking(u1.ID, r3.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	seedBooking(u1.ID, r4.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	seedBooking(u2.ID, r1.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	seedBooking(u2.ID, r2.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	seedBooking(u2.ID, r3.ID, time.Now(), time.Now().AddDate(0, 0, 1))
}

func seedBooking(userID, roomID primitive.ObjectID, startDate time.Time, endDate time.Time) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: startDate,
		TillDate: endDate,
	}
	insertedBooking, err := bookingStore.InsertBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
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
	bookingStore = db.NewMongoBookingStore(client)
}
