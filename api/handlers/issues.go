package handlers

import (
	"fmt"
	"path/filepath"
	"unibox/db"
	"unibox/models"
	"unibox/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Issue struct {
	DB *pgxpool.Pool
}

func (h *Issue) Create(c fiber.Ctx) error {
	role := c.Locals("role").(string)
	issuer := c.Locals("user_id").(string)

	if role != "user" {
		return c.Status(401).JSON(fiber.Map{"error": "You are not authorized to access this resource"})
	}

	user, err := db.GetUserById(h.DB, c.Context(), issuer)

	if err != nil {
		return c.SendStatus(500)
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	id := uuid.New().String()
	img_url := "null"

	file, err := c.FormFile("img")
	if err == nil {
		const maxSize = 3 * 1024 * 1024

		if file.Size > maxSize {
			return c.Status(400).SendString("File too large (max 3MB)")
		}

		img_url = fmt.Sprintf("%s%s", id, filepath.Ext(file.Filename))

		err = c.SaveFile(file, "./uploads/"+img_url)
		if err != nil {
			return c.Status(500).SendString("Failed to save image !")
		}

	}

	dept := utils.RouteComplain(user, title, description)

	issue := models.Issue{
		Id:     id,
		Issuer: issuer,
		Title:  title,
		Desc:   description,
		Dept:   dept,
		Img:    img_url,
		Status: "pending",
	}

	err = db.CreateIssue(h.DB, c.Context(), issue)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func (h *Issue) Get(c fiber.Ctx) error {
	role := c.Locals("role").(string)
	user_id := c.Locals("user_id").(string)

	if role == "user" {
		issues, err := db.GetIssuesUsers(h.DB, c.Context(), user_id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": issues})
	}

	if role == "admin" {
		admin, err := db.GetAdminById(h.DB, c.Context(), user_id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		issues, err := db.GetIssuesDept(h.DB, c.Context(), admin.Department)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"data": issues})

	}

	return c.SendStatus(404)
}

func (h *Issue) Update(c fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" {
		return c.SendStatus(401)
	}

	body := new(struct {
		Id      string `json:"id"`
		Action  string `json:"action"`
		Payload string `json:"payload"`
	})

	if err := c.Bind().All(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// adminId := c.Locals("user_id").(string)

	// admin, err := db.GetAdminById(h.DB, c.Context(), adminId)

	// if err != nil {
	// 	return c.SendStatus(500)
	// }

	// if admin.Department ==

	return c.SendStatus(404)
}
