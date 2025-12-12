package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/controller/http/v1/request"
	//"github.com/evrone/go-clean-template/internal/entity"
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

// Ensure this matches how you register handlers in router.go
// For convenience we expose a constructor used by the router file above.
func NewCommentHandler(u usecase.Comment, l logger.Interface) *CommentHandler {
	return &CommentHandler{U: u, L: l, V: validator.New(validator.WithRequiredStructEnabled())}
}

// helper: uniform error responses (you can replace with your project's response helpers)
/*func errorResponse(ctx *fiber.Ctx, status int, msg string) error {
	return ctx.Status(status).JSON(map[string]string{"error": msg})
}*/

// --- Handlers -----------------------------------------------------------------

// list: GET /v1/comments?entity_id=...&page=...&limit=...&sort=asc|desc
// Пояснения в комментариях внутри.
func (h *CommentHandler) list(ctx *fiber.Ctx) error {
	// Парсим query-параметры: entity_id обязателен для фильтра
	entityIDStr := ctx.Query("entity_id")
	if entityIDStr == "" {
		return errorResponse(ctx, http.StatusBadRequest, "entity_id is required")
	}
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil || entityID <= 0 {
		return errorResponse(ctx, http.StatusBadRequest, "invalid entity_id")
	}

	// Пагинация: разумные дефолты
	pageStr := ctx.Query("page", "1")
	limitStr := ctx.Query("limit", "20")
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Сортировка (по дате)
	sort := ctx.Query("sort", "desc")
	sortAsc := false
	if sort == "asc" {
		sortAsc = true
	}

	// Вызов usecase — вся логика фильтрации/пагинации в бизнес-слое
	comments, err := h.U.ListComments(ctx.UserContext(), entityID, page, limit, sortAsc)
	if err != nil {
		h.L.Error(err, "http - v1 - comment list")
		return errorResponse(ctx, http.StatusInternalServerError, "cannot fetch comments")
	}

	return ctx.Status(http.StatusOK).JSON(comments)
}

// get: GET /v1/comments/:id
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
// С тело запроса валидируется через DTO (request.CreateComment)
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
// Хендлер парсит id и тело, валидация и вызов usecase.UpdateComment
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

// -----------------------------------------------------------------------------
// Notes:
// 1) В примере выше предполагается, что usecase.Comment содержит методы:
//    CreateComment(ctx context.Context, userID, entityID int64, text string) (*entity.Comment, error)
//    UpdateComment(ctx context.Context, id, userID int64, text string) error
//    ListComments(ctx context.Context, entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error)
//    GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
// 2) В handler'е мы не реализуем бизнес-логику — только адаптация.
// 3) Маппинг ошибок из usecase в http статусы стоит централизовать: например, через пакет errors
//    с предопределёнными переменными (ErrNotFound, ErrForbidden, ErrValidation) и функцией MapToHTTP.
// 4) Реализация Delete подразумевает отдельный метод usecase.DeleteComment(ctx, id, userID).
//    В коде выше показан placeholder; подмени на реальный вызов.
// 5) Swagger аннотации можно добавить над каждым handler методом аналогично translation example.
