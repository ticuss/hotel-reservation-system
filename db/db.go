package db

var MongoDBNameEnvName = "MONGO_DB_NAME"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Pagination struct {
	Limit int64
	Page  int64
}
