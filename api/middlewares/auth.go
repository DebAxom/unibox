package middlewares

import (
	"errors"
	"os"
	"strings"
	"unibox/db"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Claims struct {
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	jwt.RegisteredClaims
}

func Auth(DB *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid auth format"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid claims"})
		}

		user_id := ""
		role := claims.Role

		if role == "user" {
			user, err := db.GetUserById(DB, c.Context(), claims.UserID)
			if err != nil {
				return c.Status(401).JSON(fiber.Map{"error": "User not found"})
			}

			if user.TokenVersion != claims.TokenVersion {
				return c.Status(401).JSON(fiber.Map{"error": "Session invalidated"})
			}

			user_id = user.Id
		}

		if role == "admin" {
			admin, err := db.GetAdminById(DB, c.Context(), claims.UserID)
			if err != nil {
				return c.Status(401).JSON(fiber.Map{"error": "User not found"})
			}

			if admin.TokenVersion != claims.TokenVersion {
				return c.Status(401).JSON(fiber.Map{"error": "Session invalidated"})
			}

			user_id = admin.Id
		}

		c.Locals("user_id", user_id)
		c.Locals("role", role)

		return c.Next()
	}
}
