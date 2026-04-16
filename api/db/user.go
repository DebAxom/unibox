package db

import (
	"context"
	"unibox/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(db *pgxpool.Pool, c context.Context, user models.User) error {
	query := `INSERT INTO users (id, name, email, scholar_id, password, gender, hostel) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(c, query, user.Id, user.Name, user.Email, user.ScholarID, user.Password, user.Gender, user.Hostel)
	return err
}

func DeleteUser(db *pgxpool.Pool, c context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(c, query, id)
	return err
}

func GetUserById(db *pgxpool.Pool, c context.Context, id string) (models.User, error) {
	query := `SELECT id, name, email, scholar_id, password, gender, hostel, token_version FROM users WHERE id = $1`
	var user models.User

	err := db.QueryRow(c, query, id).Scan(&user.Id, &user.Name, &user.Email, &user.ScholarID, &user.Password, &user.Gender, &user.Hostel, &user.TokenVersion)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *pgxpool.Pool, c context.Context, email string) (models.User, error) {
	query := `SELECT id, name, email, scholar_id, password, gender, hostel, token_version FROM users WHERE email = $1`
	var user models.User

	err := db.QueryRow(c, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.ScholarID, &user.Password, &user.Gender, &user.Hostel, &user.TokenVersion)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func IncrementTokenVersionUser(db *pgxpool.Pool, c context.Context, id string) error {
	query := `UPDATE users SET token_version = token_version + 1 WHERE id = $1`
	_, err := db.Exec(c, query, id)
	return err
}
