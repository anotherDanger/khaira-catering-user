package middleware

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func MyMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	splitHeader := strings.SplitN(authHeader, " ", 2)
	if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenString := splitHeader[1]
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	username, ok := claims["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims username"})
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid expiration claim"})
	}

	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Expired token"})
	}

	c.Locals("token", tokenString)
	c.Locals("username", username)

	return c.Next()
}
