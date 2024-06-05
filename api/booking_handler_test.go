package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/ticuss/hotel-reservation-system/db/fixtures"
)

func TestGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	user := fixtures.AddUser(tdb.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(tdb.Store, "Bar Hotel", "Asnieres", 4, nil)
	room := fixtures.AddRoom(tdb.Store, true, "small", 99.9, hotel.ID)
	from := time.Now()
	till := time.Now().AddDate(0, 0, 5)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
	fmt.Println("Booking", booking)
}
