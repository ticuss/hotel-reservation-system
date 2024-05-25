package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type RoomType int

const (
	_ RoomType = iota // Skip 0
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	Size    string             "bson:\"size\" json:\"size\""
	Price   float64            `bson:"price" json:"price"`
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
	Seaside bool               `bson:"seaside" json:"seaside"`
}
