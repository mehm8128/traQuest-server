package model

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	Name            string      `json:"name" db:"name"`
	CompletedQuests []uuid.UUID `json:"completedQuests"`
	Score           int         `json:"score" db:"score"`
}

func GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
