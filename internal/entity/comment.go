package entity

import (
	"errors"
	"strings"
	"time"
)

type Comment struct {
	ID        int64
	EntityID  int64
	UserID    int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrEmptyCommentText = errors.New("text is empty")
	ErrNotAuthor        = errors.New("not author")
)

func NewComment(entityID, userID int64, text string) (*Comment, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrEmptyCommentText
	}

	now := time.Now().UTC()

	return &Comment{
		EntityID:  entityID,
		UserID:    userID,
		Text:      text,
		CreatedAt: now,
	}, nil
}

func (c *Comment) UpdateComment(userID int64, text string) error {
	if userID != c.UserID {
		return ErrNotAuthor
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return ErrEmptyCommentText
	}
	c.Text = text
	c.UpdatedAt = time.Now().UTC()
	return nil
}
