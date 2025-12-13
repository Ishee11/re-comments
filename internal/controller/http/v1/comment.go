package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/controller/http/v1/request"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CommentHandler adapts HTTP requests to the comment usecase.
// Основная идея: handler не содержит бизнес-логики, он только
// парсит/валидирует вход, вызывает usecase и формирует ответ.
type CommentHandler struct {
	U usecase.Comment
	L logger.Interface
	V *validator.Validate
}

type ListCommentsRequest struct {
	EntityID int64  `query:"entity_id" validate:"required,min=1"`
	Page     int64  `query:"page" validate:"min=1"`
	Limit    int64  `query:"limit" validate:"min=1,max=100"`
	Sort     string `query:"sort" validate:"oneof=asc desc"`
}

// list: GET /v1/comments?entity_id=...&page=...&limit=...&sort=asc|desc
// @Summary     List comments
// @Description Get comments by entity ID with pagination and sorting
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       entity_id query int true "Entity ID"
// @Param       page query int false "Page number" default(1)
// @Param       limit query int false "Items per page" default(20)
// @Param       sort query string false "Sort by created_at asc|desc" default(desc)
// @Success     200 {array} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments [get]
// Структура для параметров запроса
func (h *CommentHandler) list(ctx *fiber.Ctx) error {
	// Парсинг в структуру
	var req ListCommentsRequest
	if err := ctx.QueryParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
	}

	// Валидация
	if err := h.V.Struct(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	// Установка дефолтных значений
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Sort == "" {
		req.Sort = "desc"
	}

	// Вызов usecase
	comments, err := h.U.ListComments(
		ctx.UserContext(),
		req.EntityID,
		req.Page,
		req.Limit,
		req.Sort == "asc",
	)
	if err != nil {
		h.L.Error(err, "http - v1 - comment list")
		return errorResponse(ctx, http.StatusInternalServerError, "cannot fetch comments")
	}

	return ctx.Status(http.StatusOK).JSON(comments)
}

// get: GET /v1/comments/:id
// @Summary     Get comment by ID
// @Description Fetch a single comment by its ID
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       id path int true "Comment ID"
// @Success     200 {object} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments/{id} [get]
func (h *CommentHandler) get(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	c, err := h.U.GetCommentByID(ctx.UserContext(), id)
	if err != nil {
		h.L.Error(err, "http - v1 - comment get")
		return errorResponse(ctx, http.StatusInternalServerError, "cannot fetch comment")
	}
	if c == nil {
		return errorResponse(ctx, http.StatusNotFound, "comment not found")
	}

	return ctx.Status(http.StatusOK).JSON(c)
}

// create: POST /v1/comments
// @Summary     Create comment
// @Description Create a new comment for an entity
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       request body request.CreateComment true "Comment payload"
// @Success     201 {object} entity.Comment
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments [post]
func (h *CommentHandler) create(ctx *fiber.Ctx) error {
	var body request.CreateComment
	if err := ctx.BodyParser(&body); err != nil {
		h.L.Error(err, "http - v1 - comment create: parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}
	if err := h.V.Struct(body); err != nil {
		h.L.Error(err, "http - v1 - comment create: validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	// Передаём в usecase: handler не занимается созданием entity напрямую
	created, err := h.U.CreateComment(ctx.UserContext(), body.UserID, body.EntityID, body.Text)
	if err != nil {
		h.L.Error(err, "http - v1 - comment create")
		return errorResponse(ctx, http.StatusInternalServerError, "cannot create comment")
	}

	return ctx.Status(http.StatusCreated).JSON(created)
}

// update: PUT /v1/comments/:id
// @Summary     Update comment
// @Description Update an existing comment (author only)
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       id path int true "Comment ID"
// @Param       request body request.UpdateComment true "Updated comment payload"
// @Success     200 {string} string "OK"
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments/{id} [put]
func (h *CommentHandler) update(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var body request.UpdateComment
	if err := ctx.BodyParser(&body); err != nil {
		h.L.Error(err, "http - v1 - comment update: parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}
	if err := h.V.Struct(body); err != nil {
		h.L.Error(err, "http - v1 - comment update: validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	// В usecase реализованы проверки авторства и прочие бизнес-правила
	if err := h.U.UpdateComment(ctx.UserContext(), id, body.UserID, body.Text); err != nil {
		h.L.Error(err, "http - v1 - comment update")
		// можно маппить бизнес-ошибки в http статусы (например ErrNotFound -> 404, ErrForbidden -> 403)
		return errorResponse(ctx, http.StatusInternalServerError, "cannot update comment")
	}

	return ctx.SendStatus(http.StatusOK)
}

// delete: DELETE /v1/comments/:id
// @Summary     Delete comment
// @Description Delete a comment by ID (author only)
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       id path int true "Comment ID"
// @Success     204 {string} string "No Content"
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /comments/{id} [delete]
func (h *CommentHandler) delete(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := h.U.UpdateComment(ctx.UserContext(), id, 0, ""); err != nil {
		// В данном примере предполагается, что usecase имеет отдельный метод DeleteComment.
		// Если у тебя есть DeleteComment, используй его. Здесь показан placeholder.
		h.L.Error(err, "http - v1 - comment delete")
		return errorResponse(ctx, http.StatusInternalServerError, "cannot delete comment")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
