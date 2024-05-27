package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ticuss/hotel-reservation-system/db"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	fmt.Println(params)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("Invalid Credentials")
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		return fmt.Errorf("Invalid credentials")
	}
	fmt.Println(user)
	return nil
}
