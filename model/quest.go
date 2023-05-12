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
	CompletedUsers []uuid.UUID `json:"completedUsers"`
}

type TagQuest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TagID     uuid.UUID `json:"tagID" db:"tag_id"`
	QuestID   uuid.UUID `json:"questID" db:"quest_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func GetQuests(ctx context.Context, userID uuid.UUID) ([]*Quest, error) {
	var quests []*Quest
	// todo: completed怪しい
	err := db.SelectContext(ctx, &quests, "SELECT quests.id, quests.number, quests.title, quests.description, quests.level, quests.created_at, quests.updated_at, users_quests.id as completed FROM quests LEFT JOIN users_quests ON quests.id = users_quests.quest_id WHERE quests.approved = true ORDER BY number")
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

func GetUnapprovedQuests(ctx context.Context) ([]*Quest, error) {
	var quests []*Quest
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

func CompleteQuest(ctx context.Context, questID, userID uuid.UUID) error {
	createdAt := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO users_quests VALUES (?, ?, ?)", questID, userID, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func CreateQuest(ctx context.Context, title string, description string, level int, tags []uuid.UUID) (*QuestDetail, error) {
	ID := uuid.New()
	createdAt := time.Now()

	var count int
	err := db.GetContext(ctx, &count, "SELECT number FROM quests ORDER BY number DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(ctx, "INSERT INTO quests (id, number, title, description, level, approved, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", ID, count, title, description, level, false, createdAt, createdAt)
	if err != nil {
		return nil, err
	}
	tagsQuests := make([]TagQuest, len(tags))
	for i := range tags {
		tagsQuests[i] = TagQuest{
			ID:        uuid.New(),
			TagID:     tags[i],
			QuestID:   ID,
			CreatedAt: createdAt,
		}
	}
	_, err = db.NamedExec("INSERT INTO tags_quests (id, tag_id, quest_id, created_at) VALUES (:id, :tag_id, quest_id, created_at)", tagsQuests)
	if err != nil {
		return nil, err
	}

	quest, err := GetQuest(ctx, ID)
	if err != nil {
		return nil, err
	}

	return quest, nil
}

func ApproveQuest(ctx context.Context, id uuid.UUID) (*QuestDetail, error) {
	quest, err := GetQuest(ctx, id)
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
	// todo: Stags_questsの更新
	// 全部消してからあるやつを再INSERT

	quest, err := GetQuest(ctx, id)
	if err != nil {
		return nil, err
	}

	return quest, nil
}
