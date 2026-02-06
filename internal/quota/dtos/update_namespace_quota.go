package dtos

type UpdateNamespaceQuotaRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []Resource `json:"resources"`
}
