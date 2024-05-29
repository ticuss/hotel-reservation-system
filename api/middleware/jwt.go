package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorized")
	}
	claims, err := validateToken(token[0])
	if err != nil {
		return err
	}
	// Check token expiration
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)
	if time.Now().After(time.Unix(expires, 0)) {
		return fmt.Errorf("token expired")
	}
	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
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
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	fmt.Println("token :", token)
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("claims :", claims, ok)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
