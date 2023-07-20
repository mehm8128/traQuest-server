package model

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID              string      `json:"id" db:"id"`
	CompletedQuests []uuid.UUID `json:"completedQuests"`
	Score           int         `json:"score" db:"score"`
}

func GetUser(ctx context.Context, id string) (*User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, score FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	err = db.SelectContext(ctx, &user.CompletedQuests, "SELECT quest_id FROM users_quests WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	if user.CompletedQuests == nil {
		user.CompletedQuests = []uuid.UUID{}
	}

	return &user, nil
}

func IsUserExist(ctx context.Context, id string) (bool, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func CreateUser(ctx context.Context, id string) (*User, error) {
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, score) VALUES (?, ?)", id, 0)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
