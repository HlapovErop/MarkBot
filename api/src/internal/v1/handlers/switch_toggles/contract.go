package switch_toggles

type Body struct {
	Toggles map[string]interface{} `json:"toggles" binding:"required"`
}
