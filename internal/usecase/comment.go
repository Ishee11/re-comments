package usecase

import "github.com/evrone/go-clean-template/internal/entity"

//go:generate mockgen -source=../usecase/comment.go -destination=../usecase/mock/mock_comment_usecase.go -package=mock
type Comment interface {
	CreateComment(userID, entityID int64, text string) (*entity.Comment, error)
	UpdateComment(id, userID int64, text string) error
	ListComments(entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(id int64) (*entity.Comment, error)
}
