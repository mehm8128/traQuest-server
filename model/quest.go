package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Quest struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Number         int       `json:"number" db:"number"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	Level          int       `json:"level" db:"level"`
	Tags           []string  `json:"tags" db:"tags"`
	CompletedCount int       `json:"completedCount" db:"completed_count"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	Completed      bool      `json:"completed"`
}

type QuestDetail struct {
	Quest
	CompletedUsers []uuid.UUID `json:"completedUsers"`
}

func GetQuests(ctx context.Context) ([]*Quest, error) {
	var quests []*Quest
	// todo: userのcomplete状態を繋げる
	err := db.SelectContext(ctx, &quests, "SELECT * FROM quests ORDER BY number")
	if err != nil {
		return nil, err
	}
	return quests, nil
}
