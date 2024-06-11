package db

const (
	DBNAME     = "hotel-reservation-system"
	DBURI      = "mongodb://localhost:27017"
	TESTDBNAME = "hotel-reservation-test"
)

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
