package entity

type ScaffoldRequest struct {
	Template    string            `json:"template" binding:"required"`
	ProjectID   string            `json:"projectId" binding:"required"`
	Environment string            `json:"environment" binding:"required"`
	Variables   map[string]string `json:"variables" binding:"required"`
}

type ScaffoldRequests []ScaffoldRequest
