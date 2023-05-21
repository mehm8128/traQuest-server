package model

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	Name            string      `json:"name" db:"name"`
	CompletedQuests []uuid.UUID `json:"completedQuests"`
	Score           int         `json:"score" db:"score"`
}

func GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, name, score FROM users WHERE id = ?", id)
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

func IsUserExist(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func CreateUser(ctx context.Context, id uuid.UUID, name string) (*User, error) {
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, name, score) VALUES (?, ?, ?)", id, name, 0)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
