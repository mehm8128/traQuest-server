package model

import "github.com/google/uuid"

type User struct {
	Name            string      `json:"name" db:"name"`
	CompletedQuests []uuid.UUID `json:"completedQuests"`
	Score           int         `json:"score" db:"score"`
}
