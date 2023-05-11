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
	Tags        []*Tag    `json:"tags"`
	Completed   bool      `json:"completed"`
}

type QuestDetail struct {
	Quest
	CompletedUsers []*uuid.UUID `json:"completedUsers"`
}

type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func GetQuests(ctx context.Context, userID uuid.UUID) ([]*Quest, error) {
	var quests []*Quest
	// todo: completed怪しい
	err := db.SelectContext(ctx, &quests, "SELECT quests.id, quests.number, quests.title, quests.description, quests.level, quests.created_at, quests.updated_at, users_quests.id as completed FROM quests LEFT JOIN users_quests ON quests.id = users_quests.quest_id ORDER BY number")
	if err != nil {
		return nil, err
	}
	//todo: n+1
	for _, quest := range quests {
		tags, err := GetTagsByQuestID(ctx, quest.ID)
		if err != nil {
			return nil, err
		}
		quest.Tags = tags
	}

	return quests, nil
}

func GetQuest(ctx context.Context, id uuid.UUID) (*QuestDetail, error) {
	var quest QuestDetail
	err := db.GetContext(ctx, &quest, "SELECT FROM quests WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	var count int
	err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users_quests WHERE quest_id = ?", id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		quest.Completed = true
	} else {
		quest.Completed = false
	}
	tags, err := GetTagsByQuestID(ctx, id)
	if err != nil {
		return nil, err
	}
	quest.Tags = tags
	err = db.SelectContext(ctx, &quest.CompletedUsers, "SELECT user_id FROM users_quests WHERE quest_id = ?", id)
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
