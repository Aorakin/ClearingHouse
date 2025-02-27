package dto

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
