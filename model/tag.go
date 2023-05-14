package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func GetTags(ctx context.Context) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	err := db.SelectContext(ctx, &tags, "SELECT * from tags ORDER BY name")
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func GetTagsByQuestID(ctx context.Context, id uuid.UUID) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	err := db.SelectContext(ctx, &tags, "SELECT tags.id, tags.name, tags.created_at FROM tags_quests JOIN tags ON tags.id = tags_quests.tag_id WHERE quest_id = ?", id)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func CreateTag(ctx context.Context, tags []string) ([]*Tag, error) {
	newTags := make([]*Tag, 0)
	for _, tag := range tags {
		ID := uuid.New()
		createdAt := time.Now()
		newTags = append(newTags, &Tag{
			ID:        ID,
			Name:      tag,
			CreatedAt: createdAt,
		})
	}

	_, err := db.NamedExecContext(ctx, "INSERT INTO tags (id, name, created_at) VALUES (:id, :name, :created_at)", newTags)
	if err != nil {
		return nil, err
	}

	return newTags, nil
}
