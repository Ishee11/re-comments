// Package repo содержит интерфейсы репозиториев для работы с внешними источниками данных.
// Репозитории отвечают за сохранение и получение сущностей, не реализуя бизнес-логику.
// В чистой архитектуре репозитории — это слой “outer layer”.
//
// Интерфейсы определяют контракты для работы с сущностями,
// чтобы UseCase мог их использовать, не зная деталей реализации.
package repo

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

// CommentRepository определяет методы для работы с комментариями.
// Этот интерфейс используется UseCase для сохранения, изменения и получения комментариев.
// Реализация может быть любая: база данных, in-memory, веб-API и т.д.
//
//go:generate mockgen -source=../repo/comment.go -destination=../usecase/mock/mock_comment_repo.go -package=mock
type CommentRepository interface {
	CreateComment(ctx context.Context, c *entity.Comment) error
	UpdateComment(ctx context.Context, c *entity.Comment) error
	ListCommentByEntity(ctx context.Context, entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
}
