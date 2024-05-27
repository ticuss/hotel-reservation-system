package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--JWT auth")
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorized")
	}
	if err := parseToken(token[0]); err != nil {
		return err
	}
	fmt.Println("Token:", token)

	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Print("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Print("failed to parse jwt token", err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims)
	}
	return fmt.Errorf("Unauthorized")
}
