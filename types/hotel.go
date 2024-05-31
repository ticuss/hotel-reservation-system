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

type Room struct {
	Size    string             "bson:\"size\" json:\"size\""
	Price   float64            `bson:"price" json:"price"`
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
	Seaside bool               `bson:"seaside" json:"seaside"`
}
