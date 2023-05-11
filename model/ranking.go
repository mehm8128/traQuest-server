package model

import (
	"github.com/google/uuid"
)

type Rank struct {
	ID     uuid.UUID `json:"id"`
	Rank   string    `json:"rank"`
	UserID string    `json:"userId"`
	Score  int       `json:"score"`
}
