package persistent

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

// CommentRepo -.
type CommentRepo struct {
	*postgres.Postgres
}

// NewCommentRepo -.
func NewCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

func (r *CommentRepo) CreateComment(ctx context.Context, c *entity.Comment) error {
	sql, args, err := r.Builder.
		Insert("comments").
		Columns("entity_id, user_id, text").
		Values(c.EntityID, c.UserID, c.Text).
		ToSql()
	if err != nil {
		return fmt.Errorf("CommentRepo - CreateComment - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CommentRepo - CreateComment - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	var comment entity.Comment
	sql, args, err := r.Builder.
		Select("id, entity_id, user_id, text, created_at").
		From("comments").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetCommentByID - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&comment.ID, &comment.EntityID, &comment.UserID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // можно маппить в 404 в usecase/handler
		}
		return nil, fmt.Errorf("CommentRepo - GetCommentByID - row.Scan: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, c *entity.Comment) error {
	sql, args, err := r.Builder.
		Update("comments").
		Set("text", c.Text).
		Where(sq.Eq{"id": c.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("CommentRepo - UpdateComment - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CommentRepo - UpdateComment - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CommentRepo) ListCommentByEntity(ctx context.Context, entityID int64,
	page int64, limit int64, sortAsc bool,
) ([]*entity.Comment, error) {
	// защита от дурака
	if page < 1 {
		page = 1
	}
	const (
		defaultLimit = 20
		maxLimit     = 100
	)
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	limit64 := uint64(limit)
	offset := uint64((page - 1) * limit)

	// направление сортировки
	dir := "DESC"
	if sortAsc {
		dir = "ASC"
	}

	sql, args, err := r.Builder.
		Select("id", "entity_id", "user_id", "text", "created_at", "updated_at").
		From("comments").
		Where("entity_id = ?", entityID).
		OrderBy(fmt.Sprintf("created_at %s", dir)).
		Limit(limit64).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - ListCommentByEntity - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - ListCommentByEntity - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	comments := make([]*entity.Comment, 0, limit)
	for rows.Next() {
		c := &entity.Comment{}
		// предполагаемые поля в entity.Comment: ID, EntityID, UserID, Text, CreatedAt, UpdatedAt
		if err := rows.Scan(&c.ID, &c.EntityID, &c.UserID, &c.Text, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("CommentRepo - ListCommentByEntity - rows.Scan: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CommentRepo - ListCommentByEntity - rows.Err: %w", err)
	}

	return comments, nil
}
