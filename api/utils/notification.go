package utils

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Notify(db *pgxpool.Pool, ctx context.Context, title, message, userID, dept, issueID string) {

	db.Exec(ctx, `INSERT INTO notifications (user_id, issue_id, dept, title, message) VALUES ($1, $2, $3, $4, $5)`, userID, issueID, dept, title, message)
}
