package model

type Rank struct {
	Rank     string `json:"rank"`
	UserName string `json:"userName"`
	Score    int    `json:"score"`
}

func GetRanking() ([]*Rank, error) {
	var ranking []*Rank
	err := db.Select(&ranking, "SELECT name, score, RANK() OVER(ORDER BY score DESC) AS rank FROM users LIMIT 20")
	if err != nil {
		return nil, err
	}
	return ranking, nil
}
