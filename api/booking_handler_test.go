package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/api/middleware"
	"github.com/ticuss/hotel-reservation-system/db/fixtures"
	"github.com/ticuss/hotel-reservation-system/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(tdb.Store, "nonauth", "nonauth", false)
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "Bar Hotel", "Asnieres", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, true, "small", 99.9, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		bookingHandler = NewBookingHandler(tdb.Store)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.User))
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Non 200 response code received %d", resp.StatusCode)
	}

	var bookingResp *types.Booking
	fmt.Println(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected booking id %s but got %s", booking.ID, bookingResp.ID)
	}
	// non auth user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected a non 200 response code received %d", resp.StatusCode)
	}
}

func TestAdminGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "Bar Hotel", "Asnieres", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, true, "small", 99.9, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	fmt.Println(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Non 200 response code received %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("Expected 1 booking, got %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected booking id %s but got %s", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected booking id %s but got %s", booking.ID, have.ID)
	}
	// Test Not Admin
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	fmt.Println(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected a non 200 response code received %d", resp.StatusCode)
	}
}
