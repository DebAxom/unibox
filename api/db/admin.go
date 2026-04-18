package db

import (
	"context"
	"unibox/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateAdmin(db *pgxpool.Pool, c context.Context, admin models.Admin) error {
	query := `INSERT INTO admins (id, name, username, password, department) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(c, query, admin.Id, admin.Name, admin.Password, admin.Department)
	return err
}

func GetAdminById(db *pgxpool.Pool, c context.Context, id string) (models.Admin, error) {
	query := `SELECT id, username, password, dept, token_version FROM admins WHERE id = $1`

	var admin models.Admin

	err := db.QueryRow(c, query, id).Scan(&admin.Id, &admin.Username, &admin.Password, &admin.Department, &admin.TokenVersion)

	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

func GetAdminByUsername(db *pgxpool.Pool, c context.Context, username string) (models.Admin, error) {

	query := `SELECT id, username, password, dept, token_version FROM admins WHERE username = $1`

	var admin models.Admin

	err := db.QueryRow(c, query, username).Scan(&admin.Id, &admin.Username, &admin.Password, &admin.Department, &admin.TokenVersion)

	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

func IncrementTokenVersionAdmin(db *pgxpool.Pool, c context.Context, id string) error {
	query := `UPDATE admins SET token_version = token_version + 1 WHERE id = $1`
	_, err := db.Exec(c, query, id)
	return err
}
