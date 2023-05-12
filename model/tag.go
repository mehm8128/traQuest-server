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

func CreateTag(ctx context.Context, name string) (*Tag, error) {
	ID := uuid.New()
	createdAt := time.Now()

	_, err := db.ExecContext(ctx, "INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)", ID, name, createdAt)
	if err != nil {
		return nil, err
	}

	tag := &Tag{
		ID:        ID,
		Name:      name,
		CreatedAt: createdAt,
	}

	return tag, nil
}
