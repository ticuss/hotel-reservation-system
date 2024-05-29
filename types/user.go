package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 8
)

type UserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params UserParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLastNameLength)
	}
	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLength)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = "invalid email"
	}
	return errors
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

type User struct {
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}

func NewUserFromParams(params UserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
