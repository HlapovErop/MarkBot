package donate

type Body struct {
	Points int64 `json:"points" binding:"required"`
	UserID uint  `json:"user_id" binding:"required"`
}
