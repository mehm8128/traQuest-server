package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UnapprovedQuest struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Number      int       `json:"number" db:"number"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Level       int       `json:"level" db:"level"`
	Approved    bool      `json:"approved" db:"approved"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Tags        []*Tag    `json:"tags"`
}

type Quest struct {
	UnapprovedQuest
	Completed bool `json:"completed"`
}

type QuestDetail struct {
	Quest
	CompletedUsers []string `json:"completedUsers"`
}

type TagQuest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TagID     uuid.UUID `json:"tagID" db:"tag_id"`
	QuestID   uuid.UUID `json:"questID" db:"quest_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func GetQuests(ctx context.Context, userID string) ([]*Quest, error) {
	quests := make([]*Quest, 0)
	err := db.SelectContext(ctx, &quests, "SELECT quests.id, quests.number, quests.title, quests.description, quests.level, quests.created_at, quests.updated_at FROM quests WHERE quests.approved = true ORDER BY number")
	if err != nil {
		return nil, err
	}

	//todo: n+1
	for _, quest := range quests {
		err := db.GetContext(ctx, &quest.Completed, "SELECT EXISTS(SELECT * FROM users_quests WHERE user_id = ? AND quest_id = ?)", userID, quest.ID)
		if err != nil {
			return nil, err
		}
		tags, err := GetTagsByQuestID(ctx, quest.ID)
		if err != nil {
			return nil, err
		}
		quest.Tags = tags
	}

	return quests, nil
}

func GetUnapprovedQuests(ctx context.Context) ([]*Quest, error) {
	quests := make([]*Quest, 0)
	err := db.SelectContext(ctx, &quests, "SELECT * from quests WHERE approved = false ORDER BY number")
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

func GetQuest(ctx context.Context, id uuid.UUID, userId string) (*QuestDetail, error) {
	var quest QuestDetail
	err := db.GetContext(ctx, &quest, "SELECT * FROM quests WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	var count int
	if userId != "" {
		err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users_quests WHERE user_id = ? && quest_id = ?", userId, id)
		if err != nil {
			return nil, err
		}
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
	completedUsers := make([]string, 0)
	err = db.SelectContext(ctx, &completedUsers, "SELECT users.name FROM users_quests JOIN users ON users.id = users_quests.user_id WHERE quest_id = ?", id)
	if err != nil {
		return nil, err
	}
	quest.CompletedUsers = completedUsers

	return &quest, nil
}

func CompleteQuest(ctx context.Context, questID, userID uuid.UUID) error {
	uuid := uuid.New()
	createdAt := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO users_quests (id, quest_id, user_id, created_at) VALUES (?, ?, ?, ?)", uuid, questID, userID, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func CreateQuest(ctx context.Context, title string, description string, level int, tags []uuid.UUID) (*QuestDetail, error) {
	ID := uuid.New()
	createdAt := time.Now()

	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM quests")
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO quests (id, number, title, description, level, approved, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", ID, count+1, title, description, level, false, createdAt, createdAt)
	if err != nil {
		return nil, err
	}

	if len(tags) != 0 {
		tagsQuests := make([]TagQuest, len(tags))
		for i := range tags {
			tagsQuests[i] = TagQuest{
				ID:        uuid.New(),
				TagID:     tags[i],
				QuestID:   ID,
				CreatedAt: createdAt,
			}
		}

		_, err = db.NamedExecContext(ctx, "INSERT INTO tags_quests (id, tag_id, quest_id, created_at) VALUES (:id, :tag_id, :quest_id, :created_at)", tagsQuests)
		if err != nil {
			return nil, err
		}
	}
	quest, err := GetQuest(ctx, ID, "")
	if err != nil {
		return nil, err
	}

	return quest, nil
}

func RejectQuest(ctx context.Context, id uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM quests WHERE id =?", id)
	if err != nil {
		return err
	}

	return nil
}

func ApproveQuest(ctx context.Context, id uuid.UUID) (*QuestDetail, error) {
	quest, err := GetQuest(ctx, id, "")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, "UPDATE quests SET approved=true WHERE id =?", id)
	if err != nil {
		return nil, err
	}

	return quest, nil
}

func UpdateQuest(ctx context.Context, id uuid.UUID, title string, description string, level int, tags []uuid.UUID) (*QuestDetail, error) {
	updatedAt := time.Now()

	_, err := db.ExecContext(ctx, "UPDATE quests SET id=?, title=?, description=?, level=?, approved=?, updated_at=?) VALUES (?, ?, ?, ?, ?, ?)", id, title, description, level, false, updatedAt)
	if err != nil {
		return nil, err
	}
	// todo: tags_questsの更新
	// 全部消してからあるやつを再INSERT

	quest, err := GetQuest(ctx, id, "")
	if err != nil {
		return nil, err
	}

	return quest, nil
}
