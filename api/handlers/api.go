package handlers

import (
	"fmt"
	"time"
	"unibox/db"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Api struct {
	DB *pgxpool.Pool
}

func (h *Api) Init(c fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role == "user" {
		user, err := db.GetUserById(h.DB, c.Context(), userId)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{"id": userId, "role": role, "name": user.Name, "email": user.Email, "scholar_id": user.ScholarID, "gender": user.Gender, "hostel": user.Hostel})
	}

	if role == "admin" {
		admin, err := db.GetAdminById(h.DB, c.Context(), userId)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{"id": userId, "role": role, "username": admin.Username, "dept": admin.Department})
	}

	return c.SendStatus(404)
}

func (h *Api) NewNotifications(c fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "user" {
		return c.SendStatus(404)
	}

	var count int

	err := h.DB.QueryRow(c.Context(), `
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND read = FALSE
	`, userId).Scan(&count)

	if err != nil {
		return c.JSON(fiber.Map{"new": false})
	}

	if count > 0 {

		return c.JSON(fiber.Map{"new": true})
	}

	return c.JSON(fiber.Map{"new": false})
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	IssueID   string    `json:"issue_id"`
	Dept      string    `json:"dept"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	Timestamp time.Time `json:"timestamp"`
}

func (h *Api) AllNotifications(c fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "user" {
		return c.SendStatus(404)
	}

	rows, err := h.DB.Query(c.Context(), `
		SELECT id, user_id, issue_id, dept, title, message, read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var n Notification

		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.IssueID,
			&n.Dept,
			&n.Title,
			&n.Message,
			&n.Read,
			&n.Timestamp,
		)
		if err != nil {
			fmt.Println("scan error:", err)
			continue
		}

		notifications = append(notifications, n)
	}

	_, err = h.DB.Exec(c.Context(), `
		UPDATE notifications
		SET read = TRUE
		WHERE user_id = $1 AND read = FALSE
	`, userId)

	return c.JSON(notifications)
}
