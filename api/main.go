package main

import (
	"context"
	"log"
	"os"
	"time"
	"unibox/handlers"
	"unibox/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var db *pgxpool.Pool  // PostgreSQL
var rdb *redis.Client // Redis

func main() {
	godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connStr := os.Getenv("DATABASE_URL")

	// Initialize PostgreSQL
	var err error
	db, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database : %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed : %v\n", err)
	}

	// Initialize Redis
	redisURL := os.Getenv("REDIS_URL")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "Unibox API",
		AppName:       "Unibox",
	})

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	AuthHandler := &handlers.Auth{DB: db, RDB: rdb}
	auth := app.Group("/auth")
	auth.Post("/otp", AuthHandler.RequestOTP)
	auth.Post("/register", AuthHandler.RegisterUser)
	auth.Post("/login", AuthHandler.LogInUser)
	auth.Post("/login/admin", AuthHandler.LogInAdmin)
	auth.Post("/logout", AuthHandler.LogOut)
	auth.Post("/logoutall", AuthHandler.LogOutAll)
	auth.Post("/refresh", AuthHandler.Refresh)

	ProfileHandler := &handlers.Profile{DB: db}
	api := app.Group("/api", middlewares.Auth(db))
	api.Get("/me", ProfileHandler.Me)

	app.Listen(":5000", fiber.ListenConfig{EnablePrefork: false})

}
