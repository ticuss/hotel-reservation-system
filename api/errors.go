package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Err  string `json:"error"`
	Code int    `json:"code"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, "Internal Server Error")
	return c.Status(apiError.Code).JSON(apiError)
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) *Error {
	return &Error{Code: code, Err: err}
}

func ErrInvalidID() Error {
	return Error{Code: http.StatusBadRequest, Err: "Invalid ID"}
}

func ErrUnauthorized() Error {
	return Error{Code: http.StatusUnauthorized, Err: "Unauthorized"}
}

func ErrBadRequest() Error {
	return Error{Code: http.StatusBadRequest, Err: "Invalid json request"}
}

func ErrResourceNotFound(res string) Error {
	return Error{Code: http.StatusNotFound, Err: res + " not found"}
}
