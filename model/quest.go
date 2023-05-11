package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Quest struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Number      int       `json:"number" db:"number"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Level       int       `json:"level" db:"level"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Tags        []string  `json:"tags"`
	Completed   bool      `json:"completed"`
}

type QuestDetail struct {
	Quest
	CompletedUsers []uuid.UUID `json:"completedUsers"`
}

type Tag struct {
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
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

func GetQuest(ctx context.Context, id uuid.UUID) (*QuestDetail, error) {
	var quest QuestDetail
	// todo: userのcomplete状態を繋げる
	err := db.GetContext(ctx, &quest, "SELECT  FROM quests WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &quest, nil
}

func GetTagsByQuestID(ctx context.Context, id uuid.UUID) ([]*Tag, error) {
	var tags []*Tag
	err := db.SelectContext(ctx, &tags, "SELECT * FROM tags WHERE quest_id = ?", id)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
