package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID   primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"base_price" json:"base_price"`
	Price     float64            `bson:"price" json:"price"`
}
