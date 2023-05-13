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
	err := db.GetContext(ctx, &user, "SELECT name, score FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	err = db.SelectContext(ctx, &user.CompletedQuests, "SELECT quest_id FROM completed_quests WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
