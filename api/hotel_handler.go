package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/db"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var queryParams HotelQueryParams
	if err := c.QueryParser(&queryParams); err != nil {
		return err
	}
	hotels, err := h.hotelStore.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
