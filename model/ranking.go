package model

import "context"

type Rank struct {
	Rank     string `json:"rank" db:"rank"`
	UserName string `json:"userName" db:"name"`
	Score    int    `json:"score" db:"score"`
}

func GetRanking(ctx context.Context) ([]*Rank, error) {
	ranking := make([]*Rank, 0)
	err := db.SelectContext(ctx, &ranking, "SELECT name, score, RANK() OVER(ORDER BY score DESC) AS rank FROM users LIMIT 20")
	if err != nil {
		return nil, err
	}
	return ranking, nil
}
