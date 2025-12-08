// Package repo содержит интерфейсы репозиториев для работы с внешними источниками данных.
// Репозитории отвечают за сохранение и получение сущностей, не реализуя бизнес-логику.
// В чистой архитектуре репозитории — это слой “outer layer”.
//
// Интерфейсы определяют контракты для работы с сущностями,
// чтобы UseCase мог их использовать, не зная деталей реализации.
package repo

import (
	"github.com/evrone/go-clean-template/internal/entity"
)

// CommentRepository определяет методы для работы с комментариями.
// Этот интерфейс используется UseCase для сохранения, изменения и получения комментариев.
// Реализация может быть любая: база данных, in-memory, веб-API и т.д.
type CommentRepository interface {
	CreateComment(c *entity.Comment) error
	UpdateComment(c *entity.Comment) error
	ListCommentByEntity(entityID int64, page int64, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(id int64) (*entity.Comment, error)
}
