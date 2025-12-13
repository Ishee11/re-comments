// Package comment реализует бизнес-логику работы с комментариями.
//
// Пакет содержит UseCase для управления комментариями: создание, редактирование и получение списка.
// UseCase оперирует сущностями Comment из пакета entity и использует репозиторий для хранения данных.
// Пакет не зависит от конкретной базы данных или HTTP/AMQP, это чистый бизнес-слой.
//
// Архитектурно:
// - Сущности (entity.Comment) содержат бизнес-логику и проверки.
// - UseCase координирует работу с сущностями и репозиториями.
// - Репозитории реализуются отдельно и подставляются через интерфейсы.
package comment

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase реализует бизнес-логику работы с комментариями.
// Содержит методы для создания, изменения и получения списка комментариев.
// Использует репозиторий для сохранения и получения данных.
type UseCase struct {
	repo repo.CommentRepository
}

// New создаёт новый CommentUseCase с переданным репозиторием.
func New(repo repo.CommentRepository) *UseCase {
	return &UseCase{repo: repo}
}

// CreateComment создаёт новый комментарий для указанной сущности.
// Проверяет текст комментария через сущность Comment.
// Сохраняет комментарий через репозиторий.
// Возвращает созданный комментарий или ошибку.
func (uc *UseCase) CreateComment(ctx context.Context, userID, entityID int64, text string) (*entity.Comment, error) {
	c, err := entity.NewComment(entityID, userID, text)
	if err != nil {
		return nil, fmt.Errorf("CommentUseCase - CreateComment - entity.NewComment: %w", err)
	}
	if err = uc.repo.CreateComment(ctx, c); err != nil {
		return nil, fmt.Errorf("CommentUseCase - CreateComment - uc.repo.CreateComment: %w", err)
	}
	return c, nil
}

// UpdateComment изменяет текст существующего комментария.
// Проверяет, что пользователь является автором комментария.
// Вызывает метод UpdateCommentText сущности Comment и сохраняет изменения в репозитории.
// Возвращает ошибку, если комментарий не найден или пользователь не автор.
func (uc *UseCase) UpdateComment(ctx context.Context, id, userID int64, text string) error {
	c, err := uc.repo.GetCommentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CommentUseCase - UpdateComment - uc.repo.GetCommentByID: %w", err)
	}
	if err = c.UpdateComment(userID, text); err != nil {
		return fmt.Errorf("CommentUseCase - UpdateCommentText - c.UpdateCommentText: %w", err)
	}
	return uc.repo.UpdateComment(ctx, c)
}

// ListComments возвращает список комментариев для указанной сущности.
// Поддерживает пагинацию через page и limit, а также сортировку по дате через sortAsc.
// Обращается к репозиторию для получения данных.
// Возвращает срез комментариев или ошибку.
func (uc *UseCase) ListComments(ctx context.Context, entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error) {
	comments, err := uc.repo.ListCommentByEntity(ctx, entityID, page, limit, sortAsc)
	if err != nil {
		return nil, fmt.Errorf("CommentUseCase - ListComments - uc.repo.ListCommentByEntity: %w", err)
	}
	return comments, nil
}

func (uc *UseCase) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	return uc.repo.GetCommentByID(ctx, id)
}
