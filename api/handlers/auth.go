package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"
	"unibox/db"
	"unibox/models"
	"unibox/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID       string `json:"user_id"`
	TokenVersion int    `json:"token_version"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, token_version int, role string, duration time.Duration, secret string) (string, error) {
	claims := Claims{
		UserID:       id,
		TokenVersion: token_version,
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(9000)) // 0–8999
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", n.Int64()+1000), nil // 1000–9999
}

type Auth struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

// Users
func (h *Auth) RegisterUser(c fiber.Ctx) error {
	body := new(struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		ScholarID string `json:"scholar_id"`
		Password  string `json:"password"`
		Gender    string `json:"gender"`
		Hostel    string `json:"hostel"`
		Otp       string `json:"otp"`
	})

	if err := c.Bind().All(body); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	// Handling OTP starts here
	key := "otp:" + body.Email
	storedOTP, err := h.RDB.Get(c.Context(), key).Result()
	if err != nil {
		return c.Status(410).JSON(fiber.Map{"error": "OTP expired or not found"})
	}

	if storedOTP != body.Otp {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid OTP"})
	}

	h.RDB.Del(c.Context(), key) // Delete OTP after verification
	// Handling OTP ends here

	id := uuid.New().String() + time.Now().String()

	hashedPassword, err := HashPassword(body.Password)

	if err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	user := models.User{
		Id:        id,
		Name:      body.Name,
		Email:     body.Email,
		ScholarID: body.ScholarID,
		Password:  hashedPassword,
		Gender:    body.Gender,
		Hostel:    body.Hostel,
	}

	if err := db.CreateUser(h.DB, c.Context(), user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, _ := GenerateToken(user.Id, 0, "user", 15*time.Minute, os.Getenv("ACCESS_SECRET"))
	refreshToken, _ := GenerateToken(user.Id, 0, "user", 15*24*time.Hour, os.Getenv("REFRESH_SECRET"))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Now().Add(15 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{"access_token": accessToken})

}

func (h *Auth) RequestOTP(c fiber.Ctx) error {
	body := new(struct {
		Email string `json:"email"`
	})

	if err := c.Bind().All(body); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	key := "otp:" + body.Email
	otp, _ := GenerateOTP()

	fmt.Println(key, otp)
	err := utils.SendOTP(otp, body.Email)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Couldn't send OTP")
	}

	h.RDB.Set(c.Context(), key, otp, 7*time.Minute)

	return c.SendStatus(200)
}

func (h *Auth) LogInUser(c fiber.Ctx) error {
	body := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.Bind().All(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := db.GetUserByEmail(h.DB, c.Context(), body.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	if !CheckPassword(user.Password, body.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Wrong Password"})
	}

	accessToken, _ := GenerateToken(user.Id, user.TokenVersion, "user", 15*time.Minute, os.Getenv("ACCESS_SECRET"))
	refreshToken, _ := GenerateToken(user.Id, user.TokenVersion, "user", 15*24*time.Hour, os.Getenv("REFRESH_SECRET"))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Now().Add(15 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{"access_token": accessToken})

}

// Admin
func (h *Auth) LogInAdmin(c fiber.Ctx) error {
	body := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})

	if err := c.Bind().All(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := db.GetAdminByUsername(h.DB, c.Context(), body.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	if !CheckPassword(user.Password, body.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Wrong Password"})
	}

	accessToken, _ := GenerateToken(user.Id, user.TokenVersion, "admin", 15*time.Minute, os.Getenv("ACCESS_SECRET"))
	refreshToken, _ := GenerateToken(user.Id, user.TokenVersion, "admin", 15*24*time.Hour, os.Getenv("REFRESH_SECRET"))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Now().Add(15 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{"access_token": accessToken})

}

// Common Logout functions
func (h *Auth) LogOut(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,  // true in production
		SameSite: "None", // must match original
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.SendStatus(200)
}

func (h *Auth) LogOutAll(c fiber.Ctx) error {
	tokenStr := c.Cookies("refresh_token")

	if tokenStr == "" {
		return c.SendStatus(401)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(401)
	}

	claims := token.Claims.(*Claims)
	role := claims.Role

	if role == "user" {
		db.IncrementTokenVersionUser(h.DB, c.Context(), claims.UserID)
	}

	if role == "admin" {
		db.IncrementTokenVersionAdmin(h.DB, c.Context(), claims.UserID)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,  // true in production
		SameSite: "None", // must match original
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.SendStatus(200)
}

func (h *Auth) Refresh(c fiber.Ctx) error {
	tokenStr := c.Cookies("refresh_token")

	if tokenStr == "" {
		return c.SendStatus(401)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(401)
	}

	claims := token.Claims.(*Claims)
	role := claims.Role

	token_version := 0

	if role == "user" {
		user, err := db.GetUserById(h.DB, c.Context(), claims.UserID)
		if err != nil {
			return c.SendStatus(401)
		}

		if user.TokenVersion != claims.TokenVersion {
			return c.SendStatus(401)
		}

		token_version = user.TokenVersion
	}

	if role == "admin" {
		admin, err := db.GetAdminById(h.DB, c.Context(), claims.UserID)
		if err != nil {
			return c.SendStatus(401)
		}

		if admin.TokenVersion != claims.TokenVersion {
			return c.SendStatus(401)
		}

		token_version = admin.TokenVersion
	}

	newAccess, _ := GenerateToken(claims.UserID, token_version, role, 15*time.Minute, os.Getenv("ACCESS_SECRET"))
	refreshToken, _ := GenerateToken(claims.UserID, token_version, role, 15*24*time.Hour, os.Getenv("REFRESH_SECRET"))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Now().Add(15 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{"access_token": newAccess})
}
