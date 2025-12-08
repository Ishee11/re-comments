package usecase

import "github.com/evrone/go-clean-template/internal/entity"

type Comment interface {
	CreateComment(userID, entityID int64, text string) (*entity.Comment, error)
	UpdateComment(id, userID int64, text string) error
	ListComments(entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(id int64) (*entity.Comment, error)
}
