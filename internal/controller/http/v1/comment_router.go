package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewCommentRoutes registers comment-related endpoints under the given API v1 group.
// Register this from the top-level router like:
// apiV1Group := app.Group("/v1")
// v1.NewCommentRoutes(apiV1Group, commentUseCase, logger)
func NewCommentRoutes(apiV1Group fiber.Router, u usecase.Comment, l logger.Interface) {
	h := &CommentHandler{U: u, L: l, V: validator.New(validator.WithRequiredStructEnabled())}

	group := apiV1Group.Group("/comments")
	{
		group.Get("/", h.list)         // GET /v1/comments?entity_id=123&page=1&limit=20
		group.Get("/:id", h.get)       // GET /v1/comments/:id
		group.Post("/", h.create)      // POST /v1/comments
		group.Put("/:id", h.update)    // PUT /v1/comments/:id
		group.Delete("/:id", h.delete) // DELETE /v1/comments/:id
	}
}
