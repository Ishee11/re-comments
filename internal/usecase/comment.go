package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=../usecase/comment.go -destination=../usecase/mock/mock_comment_usecase.go -package=mock
type Comment interface {
	CreateComment(ctx context.Context, userID, entityID int64, text string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, id, userID int64, text string) error
	ListComments(ctx context.Context, entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
}
