package db

import (
	"context"
	"unibox/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateIssue(db *pgxpool.Pool, c context.Context, issue models.Issue) error {
	query := `INSERT INTO issues (id, issuer, title, description, img, status, dept) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(c, query, issue.Id, issue.Issuer, issue.Title, issue.Desc, issue.Img, issue.Status, issue.Dept)
	return err
}

func GetIssuesUsers(db *pgxpool.Pool, c context.Context, user_id string) ([]models.Issue, error) {

	query := `
		SELECT id, issuer, title, description, img, status, dept, updated_at 
		FROM issues
		WHERE issuer = $1
	`

	rows, err := db.Query(c, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var issue models.Issue

		err := rows.Scan(
			&issue.Id,
			&issue.Issuer,
			&issue.Title,
			&issue.Desc,
			&issue.Img,
			&issue.Status,
			&issue.Dept,
			&issue.Updated_at,
		)

		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

func GetIssuesDept(db *pgxpool.Pool, c context.Context, dept string) ([]models.Issue, error) {

	query := `
		SELECT id, issuer, title, description, img, status, dept, updated_at
		FROM issues
		WHERE dept = $1
	`

	rows, err := db.Query(c, query, dept)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var issue models.Issue

		err := rows.Scan(
			&issue.Id,
			&issue.Issuer,
			&issue.Title,
			&issue.Desc,
			&issue.Img,
			&issue.Status,
			&issue.Dept,
			&issue.Updated_at,
		)

		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}
