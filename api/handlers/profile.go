package handlers

import (
	"unibox/db"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	DB *pgxpool.Pool
}

func (h *Profile) Me(c fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role == "user" {
		user, err := db.GetUserById(h.DB, c.Context(), userId)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{"id": userId, "role": role, "name": user.Name, "email": user.Email, "scholar_id": user.ScholarID, "gender": user.Gender, "hostel": user.Hostel})
	}

	if role == "user" {
		admin, err := db.GetAdminById(h.DB, c.Context(), userId)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(fiber.Map{"id": userId, "role": role, "name": admin.Name, "dept": admin.Department})
	}

	return c.SendStatus(404)
}
