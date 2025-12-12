package request

type CreateComment struct {
	EntityID int64  `json:"entity_id" validate:"required,gt=0"`
	UserID   int64  `json:"user_id" validate:"required,gt=0"`
	Text     string `json:"text" validate:"required,min=1,max=2000"`
}

type UpdateComment struct {
	Text string `json:"text" validate:"required,min=1,max=2000"`
}
