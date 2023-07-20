package model

import "context"

type Rank struct {
	Rank  string `json:"rank" db:"rank"`
	ID    string `json:"id" db:"id"`
	Score int    `json:"score" db:"score"`
}

func GetRanking(ctx context.Context) ([]*Rank, error) {
	ranking := make([]*Rank, 0)
	err := db.SelectContext(ctx, &ranking, "SELECT id, score, RANK() OVER(ORDER BY score DESC) AS rank FROM users LIMIT 20")
	if err != nil {
		return nil, err
	}
	return ranking, nil
}
