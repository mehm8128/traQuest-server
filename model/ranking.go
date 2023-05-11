package model

import (
	"github.com/google/uuid"
)

type Rank struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Rank   string    `json:"rank" db:"rank"`
	UserID string    `json:"userId" db:"user_id"`
	Score  int       `json:"score" db:"score"`
}
